// Package vmess contains protocol definition, io lib for VMess.
package protocol

import (
	"encoding/binary"
	"hash/fnv"
	"io"
	"time"

	"github.com/v2ray/v2ray-core/common/alloc"
	v2crypto "github.com/v2ray/v2ray-core/common/crypto"
	"github.com/v2ray/v2ray-core/common/log"
	v2net "github.com/v2ray/v2ray-core/common/net"
	proxyerrors "github.com/v2ray/v2ray-core/proxy/common/errors"
	"github.com/v2ray/v2ray-core/proxy/vmess"
	"github.com/v2ray/v2ray-core/proxy/vmess/protocol/user"
	"github.com/v2ray/v2ray-core/transport"
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

// VMessRequest implements the request message of VMess protocol. It only contains the header of a
// request message. The data part will be handled by conection handler directly, in favor of data
// streaming.
type VMessRequest struct {
	Version        byte
	User           vmess.User
	RequestIV      []byte
	RequestKey     []byte
	ResponseHeader []byte
	Command        byte
	Address        v2net.Address
	Port           v2net.Port
}

// Destination is the final destination of this request.
func (this *VMessRequest) Destination() v2net.Destination {
	if this.Command == CmdTCP {
		return v2net.TCPDestination(this.Address, this.Port)
	} else {
		return v2net.UDPDestination(this.Address, this.Port)
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
func (this *VMessRequestReader) Read(reader io.Reader) (*VMessRequest, error) {
	buffer := alloc.NewSmallBuffer()

	nBytes, err := v2net.ReadAllBytes(reader, buffer.Value[:vmess.IDBytesLen])
	if err != nil {
		return nil, err
	}

	userObj, timeSec, valid := this.vUserSet.GetUser(buffer.Value[:nBytes])
	if !valid {
		return nil, proxyerrors.InvalidAuthentication
	}

	aesStream, err := v2crypto.NewAesDecryptionStream(userObj.ID().CmdKey(), user.Int64Hash(timeSec))
	if err != nil {
		return nil, err
	}

	decryptor := v2crypto.NewCryptionReader(aesStream, reader)

	nBytes, err = v2net.ReadAllBytes(decryptor, buffer.Value[:41])
	if err != nil {
		return nil, err
	}
	bufferLen := nBytes

	request := &VMessRequest{
		User:    userObj,
		Version: buffer.Value[0],
	}

	if request.Version != Version {
		log.Warning("Invalid protocol version %d", request.Version)
		return nil, proxyerrors.InvalidProtocolVersion
	}

	request.RequestIV = buffer.Value[1:17]       // 16 bytes
	request.RequestKey = buffer.Value[17:33]     // 16 bytes
	request.ResponseHeader = buffer.Value[33:37] // 4 bytes
	request.Command = buffer.Value[37]

	request.Port = v2net.PortFromBytes(buffer.Value[38:40])

	switch buffer.Value[40] {
	case addrTypeIPv4:
		_, err = v2net.ReadAllBytes(decryptor, buffer.Value[41:45]) // 4 bytes
		bufferLen += 4
		if err != nil {
			return nil, err
		}
		request.Address = v2net.IPAddress(buffer.Value[41:45])
	case addrTypeIPv6:
		_, err = v2net.ReadAllBytes(decryptor, buffer.Value[41:57]) // 16 bytes
		bufferLen += 16
		if err != nil {
			return nil, err
		}
		request.Address = v2net.IPAddress(buffer.Value[41:57])
	case addrTypeDomain:
		_, err = v2net.ReadAllBytes(decryptor, buffer.Value[41:42])
		if err != nil {
			return nil, err
		}
		domainLength := int(buffer.Value[41])
		_, err = v2net.ReadAllBytes(decryptor, buffer.Value[42:42+domainLength])
		if err != nil {
			return nil, err
		}
		bufferLen += 1 + domainLength
		request.Address = v2net.DomainAddress(string(buffer.Value[42 : 42+domainLength]))
	}

	_, err = v2net.ReadAllBytes(decryptor, buffer.Value[bufferLen:bufferLen+4])
	if err != nil {
		return nil, err
	}

	fnv1a := fnv.New32a()
	fnv1a.Write(buffer.Value[:bufferLen])
	actualHash := fnv1a.Sum32()
	expectedHash := binary.BigEndian.Uint32(buffer.Value[bufferLen : bufferLen+4])

	if actualHash != expectedHash {
		return nil, transport.CorruptedPacket
	}

	return request, nil
}

// ToBytes returns a VMessRequest in the form of byte array.
func (this *VMessRequest) ToBytes(idHash user.CounterHash, randomRangeInt64 user.RandomInt64InRange, buffer *alloc.Buffer) (*alloc.Buffer, error) {
	if buffer == nil {
		buffer = alloc.NewSmallBuffer().Clear()
	}

	counter := randomRangeInt64(time.Now().Unix(), 30)
	hash := idHash.Hash(this.User.ID().Bytes(), counter)

	buffer.Append(hash)

	encryptionBegin := buffer.Len()

	buffer.AppendBytes(this.Version)
	buffer.Append(this.RequestIV)
	buffer.Append(this.RequestKey)
	buffer.Append(this.ResponseHeader)
	buffer.AppendBytes(this.Command)
	buffer.Append(this.Port.Bytes())

	switch {
	case this.Address.IsIPv4():
		buffer.AppendBytes(addrTypeIPv4)
		buffer.Append(this.Address.IP())
	case this.Address.IsIPv6():
		buffer.AppendBytes(addrTypeIPv6)
		buffer.Append(this.Address.IP())
	case this.Address.IsDomain():
		buffer.AppendBytes(addrTypeDomain, byte(len(this.Address.Domain())))
		buffer.Append([]byte(this.Address.Domain()))
	}

	encryptionEnd := buffer.Len()

	fnv1a := fnv.New32a()
	fnv1a.Write(buffer.Value[encryptionBegin:encryptionEnd])

	fnvHash := fnv1a.Sum32()
	buffer.AppendBytes(byte(fnvHash>>24), byte(fnvHash>>16), byte(fnvHash>>8), byte(fnvHash))
	encryptionEnd += 4

	aesStream, err := v2crypto.NewAesEncryptionStream(this.User.ID().CmdKey(), user.Int64Hash(counter))
	if err != nil {
		return nil, err
	}
	aesStream.XORKeyStream(buffer.Value[encryptionBegin:encryptionEnd], buffer.Value[encryptionBegin:encryptionEnd])

	return buffer, nil
}
