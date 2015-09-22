// Package vmess contains protocol definition, io lib for VMess.
package protocol

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
	"errors"
	"hash/fnv"
	"io"
	"time"

	v2io "github.com/v2ray/v2ray-core/common/io"
	"github.com/v2ray/v2ray-core/common/log"
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/proxy/vmess/protocol/user"
)

const (
	addrTypeIPv4   = byte(0x01)
	addrTypeIPv6   = byte(0x03)
	addrTypeDomain = byte(0x02)

	CmdTCP = byte(0x01)
	CmdUDP = byte(0x02)

	Version = byte(0x01)

	blockSize = 16
)

var (
	ErrorInvalidUser   = errors.New("Invalid User")
	ErrorInvalidVerion = errors.New("Invalid Version")
	ErrorInvalidHash   = errors.New("Invalid Hash")
)

// VMessRequest implements the request message of VMess protocol. It only contains the header of a
// request message. The data part will be handled by conection handler directly, in favor of data
// streaming.
type VMessRequest struct {
	Version        byte
	UserId         user.ID
	RequestIV      [16]byte
	RequestKey     [16]byte
	ResponseHeader [4]byte
	Command        byte
	Address        v2net.Address
}

// Destination is the final destination of this request.
func (request *VMessRequest) Destination() v2net.Destination {
	if request.Command == CmdTCP {
		return v2net.NewTCPDestination(request.Address)
	} else {
		return v2net.NewUDPDestination(request.Address)
	}
}

// VMessRequestReader is a parser to read VMessRequest from a byte stream.
type VMessRequestReader struct {
	vUserSet user.UserSet
}

// NewVMessRequestReader creates a new VMessRequestReader with a given UserSet
func NewVMessRequestReader(vUserSet user.UserSet) *VMessRequestReader {
	return &VMessRequestReader{
		vUserSet: vUserSet,
	}
}

// Read reads a VMessRequest from a byte stream.
func (r *VMessRequestReader) Read(reader io.Reader) (*VMessRequest, error) {
	buffer := make([]byte, 256)

	nBytes, err := reader.Read(buffer[:user.IDBytesLen])
	if err != nil {
		return nil, err
	}

	log.Debug("Read user hash: %v", buffer[:nBytes])

	userId, timeSec, valid := r.vUserSet.GetUser(buffer[:nBytes])
	if !valid {
		return nil, ErrorInvalidUser
	}

	aesCipher, err := aes.NewCipher(userId.CmdKey())
	if err != nil {
		return nil, err
	}
	aesStream := cipher.NewCFBDecrypter(aesCipher, user.Int64Hash(timeSec))
	decryptor := v2io.NewCryptionReader(aesStream, reader)

	if err != nil {
		return nil, err
	}

	nBytes, err = decryptor.Read(buffer[:41])
	if err != nil {
		return nil, err
	}
	bufferLen := nBytes

	request := &VMessRequest{
		UserId:  *userId,
		Version: buffer[0],
	}

	if request.Version != Version {
		log.Error("Unknown VMess version %d", request.Version)
		return nil, ErrorInvalidVerion
	}

	copy(request.RequestIV[:], buffer[1:17])       // 16 bytes
	copy(request.RequestKey[:], buffer[17:33])     // 16 bytes
	copy(request.ResponseHeader[:], buffer[33:37]) // 4 bytes
	request.Command = buffer[37]

	port := binary.BigEndian.Uint16(buffer[38:40])

	switch buffer[40] {
	case addrTypeIPv4:
		_, err = decryptor.Read(buffer[41:45]) // 4 bytes
		bufferLen += 4
		if err != nil {
			return nil, err
		}
		request.Address = v2net.IPAddress(buffer[41:45], port)
	case addrTypeIPv6:
		_, err = decryptor.Read(buffer[41:57]) // 16 bytes
		bufferLen += 16
		if err != nil {
			return nil, err
		}
		request.Address = v2net.IPAddress(buffer[41:57], port)
	case addrTypeDomain:
		_, err = decryptor.Read(buffer[41:42])
		if err != nil {
			return nil, err
		}
		domainLength := int(buffer[41])
		_, err = decryptor.Read(buffer[42 : 42+domainLength])
		if err != nil {
			return nil, err
		}
		bufferLen += 1 + domainLength
		request.Address = v2net.DomainAddress(string(buffer[42:42+domainLength]), port)
	}

	_, err = decryptor.Read(buffer[bufferLen : bufferLen+4])
	if err != nil {
		return nil, err
	}

	fnv1a := fnv.New32a()
	fnv1a.Write(buffer[:bufferLen])
	actualHash := fnv1a.Sum32()
	expectedHash := binary.BigEndian.Uint32(buffer[bufferLen : bufferLen+4])

	if actualHash != expectedHash {
		return nil, ErrorInvalidHash
	}

	return request, nil
}

// ToBytes returns a VMessRequest in the form of byte array.
func (request *VMessRequest) ToBytes(idHash user.CounterHash, randomRangeInt64 user.RandomInt64InRange, buffer []byte) ([]byte, error) {
	if buffer == nil {
		buffer = make([]byte, 0, 300)
	}

	counter := randomRangeInt64(time.Now().UTC().Unix(), 30)
	hash := idHash.Hash(request.UserId.Bytes, counter)

	log.Debug("Writing userhash: %v", hash)
	buffer = append(buffer, hash...)

	encryptionBegin := len(buffer)

	buffer = append(buffer, request.Version)
	buffer = append(buffer, request.RequestIV[:]...)
	buffer = append(buffer, request.RequestKey[:]...)
	buffer = append(buffer, request.ResponseHeader[:]...)
	buffer = append(buffer, request.Command)
	buffer = append(buffer, request.Address.PortBytes()...)

	switch {
	case request.Address.IsIPv4():
		buffer = append(buffer, addrTypeIPv4)
		buffer = append(buffer, request.Address.IP()...)
	case request.Address.IsIPv6():
		buffer = append(buffer, addrTypeIPv6)
		buffer = append(buffer, request.Address.IP()...)
	case request.Address.IsDomain():
		buffer = append(buffer, addrTypeDomain)
		buffer = append(buffer, byte(len(request.Address.Domain())))
		buffer = append(buffer, []byte(request.Address.Domain())...)
	}

	encryptionEnd := len(buffer)

	fnv1a := fnv.New32a()
	fnv1a.Write(buffer[encryptionBegin:encryptionEnd])

	fnvHash := fnv1a.Sum32()
	buffer = append(buffer, byte(fnvHash>>24))
	buffer = append(buffer, byte(fnvHash>>16))
	buffer = append(buffer, byte(fnvHash>>8))
	buffer = append(buffer, byte(fnvHash))
	encryptionEnd += 4

	aesCipher, err := aes.NewCipher(request.UserId.CmdKey())
	if err != nil {
		return nil, err
	}
	aesStream := cipher.NewCFBEncrypter(aesCipher, user.Int64Hash(counter))
	aesStream.XORKeyStream(buffer[encryptionBegin:encryptionEnd], buffer[encryptionBegin:encryptionEnd])

	return buffer, nil
}

// VMessResponse is the header of a TCP response in VMess format.
type VMessResponse [4]byte

// NewVMessResponse creates a VMessResponse from a given VMessRequest.
func NewVMessResponse(request *VMessRequest) *VMessResponse {
	return &VMessResponse{
		request.ResponseHeader[0],
		request.ResponseHeader[1],
		request.ResponseHeader[2],
		request.ResponseHeader[3]}
}
