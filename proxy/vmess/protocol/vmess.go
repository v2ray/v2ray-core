// Package vmess contains protocol definition, io lib for VMess.
package protocol

import (
	"crypto/md5"
	"encoding/binary"
	"hash/fnv"
	"io"

	"github.com/v2ray/v2ray-core/common/alloc"
	v2crypto "github.com/v2ray/v2ray-core/common/crypto"
	"github.com/v2ray/v2ray-core/common/log"
	v2net "github.com/v2ray/v2ray-core/common/net"
	proto "github.com/v2ray/v2ray-core/common/protocol"
	"github.com/v2ray/v2ray-core/proxy"
	"github.com/v2ray/v2ray-core/transport"
)

const (
	addrTypeIPv4   = byte(0x01)
	addrTypeIPv6   = byte(0x03)
	addrTypeDomain = byte(0x02)

	CmdTCP = byte(0x01)
	CmdUDP = byte(0x02)

	Version = byte(0x01)

	OptionChunk = byte(0x01)

	blockSize = 16
)

func hashTimestamp(t proto.Timestamp) []byte {
	once := t.Bytes()
	bytes := make([]byte, 0, 32)
	bytes = append(bytes, once...)
	bytes = append(bytes, once...)
	bytes = append(bytes, once...)
	bytes = append(bytes, once...)
	return bytes
}

// VMessRequest implements the request message of VMess protocol. It only contains the header of a
// request message. The data part will be handled by connection handler directly, in favor of data
// streaming.
type VMessRequest struct {
	Version        byte
	User           *proto.User
	RequestIV      []byte
	RequestKey     []byte
	ResponseHeader byte
	Command        byte
	Option         byte
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

func (this *VMessRequest) IsChunkStream() bool {
	return (this.Option & OptionChunk) == OptionChunk
}

// VMessRequestReader is a parser to read VMessRequest from a byte stream.
type VMessRequestReader struct {
	vUserSet proto.UserValidator
}

// NewVMessRequestReader creates a new VMessRequestReader with a given UserSet
func NewVMessRequestReader(vUserSet proto.UserValidator) *VMessRequestReader {
	return &VMessRequestReader{
		vUserSet: vUserSet,
	}
}

// Read reads a VMessRequest from a byte stream.
func (this *VMessRequestReader) Read(reader io.Reader) (*VMessRequest, error) {
	buffer := alloc.NewSmallBuffer()
	defer buffer.Release()

	nBytes, err := io.ReadFull(reader, buffer.Value[:proto.IDBytesLen])
	if err != nil {
		log.Debug("VMess: Failed to read request ID (", nBytes, " bytes): ", err)
		return nil, err
	}

	userObj, timeSec, valid := this.vUserSet.Get(buffer.Value[:nBytes])
	if !valid {
		return nil, proxy.ErrorInvalidAuthentication
	}

	timestampHash := TimestampHash()
	timestampHash.Write(hashTimestamp(timeSec))
	iv := timestampHash.Sum(nil)
	aesStream, err := v2crypto.NewAesDecryptionStream(userObj.ID.CmdKey(), iv)
	if err != nil {
		log.Debug("VMess: Failed to create AES stream: ", err)
		return nil, err
	}

	decryptor := v2crypto.NewCryptionReader(aesStream, reader)

	nBytes, err = io.ReadFull(decryptor, buffer.Value[:41])
	if err != nil {
		log.Debug("VMess: Failed to read request header (", nBytes, " bytes): ", err)
		return nil, err
	}
	bufferLen := nBytes

	request := &VMessRequest{
		User:    userObj,
		Version: buffer.Value[0],
	}

	if request.Version != Version {
		log.Warning("VMess: Invalid protocol version ", request.Version)
		return nil, proxy.ErrorInvalidProtocolVersion
	}

	request.RequestIV = append([]byte(nil), buffer.Value[1:17]...)   // 16 bytes
	request.RequestKey = append([]byte(nil), buffer.Value[17:33]...) // 16 bytes
	request.ResponseHeader = buffer.Value[33]                        // 1 byte
	request.Option = buffer.Value[34]                                // 1 byte + 2 bytes reserved
	request.Command = buffer.Value[37]

	request.Port = v2net.PortFromBytes(buffer.Value[38:40])

	switch buffer.Value[40] {
	case addrTypeIPv4:
		nBytes, err = io.ReadFull(decryptor, buffer.Value[41:45]) // 4 bytes
		bufferLen += 4
		if err != nil {
			log.Debug("VMess: Failed to read target IPv4 (", nBytes, " bytes): ", err)
			return nil, err
		}
		request.Address = v2net.IPAddress(buffer.Value[41:45])
	case addrTypeIPv6:
		nBytes, err = io.ReadFull(decryptor, buffer.Value[41:57]) // 16 bytes
		bufferLen += 16
		if err != nil {
			log.Debug("VMess: Failed to read target IPv6 (", nBytes, " bytes): ", nBytes, err)
			return nil, err
		}
		request.Address = v2net.IPAddress(buffer.Value[41:57])
	case addrTypeDomain:
		nBytes, err = io.ReadFull(decryptor, buffer.Value[41:42])
		if err != nil {
			log.Debug("VMess: Failed to read target domain (", nBytes, " bytes): ", nBytes, err)
			return nil, err
		}
		domainLength := int(buffer.Value[41])
		if domainLength == 0 {
			return nil, transport.ErrorCorruptedPacket
		}
		nBytes, err = io.ReadFull(decryptor, buffer.Value[42:42+domainLength])
		if err != nil {
			log.Debug("VMess: Failed to read target domain (", nBytes, " bytes): ", nBytes, err)
			return nil, err
		}
		bufferLen += 1 + domainLength
		domainBytes := append([]byte(nil), buffer.Value[42:42+domainLength]...)
		request.Address = v2net.DomainAddress(string(domainBytes))
	}

	nBytes, err = io.ReadFull(decryptor, buffer.Value[bufferLen:bufferLen+4])
	if err != nil {
		log.Debug("VMess: Failed to read checksum (", nBytes, " bytes): ", nBytes, err)
		return nil, err
	}

	fnv1a := fnv.New32a()
	fnv1a.Write(buffer.Value[:bufferLen])
	actualHash := fnv1a.Sum32()
	expectedHash := binary.BigEndian.Uint32(buffer.Value[bufferLen : bufferLen+4])

	if actualHash != expectedHash {
		return nil, transport.ErrorCorruptedPacket
	}

	return request, nil
}

// ToBytes returns a VMessRequest in the form of byte array.
func (this *VMessRequest) ToBytes(timestampGenerator proto.TimestampGenerator, buffer *alloc.Buffer) (*alloc.Buffer, error) {
	if buffer == nil {
		buffer = alloc.NewSmallBuffer().Clear()
	}

	timestamp := timestampGenerator()
	idHash := IDHash(this.User.AnyValidID().Bytes())
	idHash.Write(timestamp.Bytes())

	hashStart := buffer.Len()
	buffer.Slice(0, hashStart+16)
	idHash.Sum(buffer.Value[hashStart:hashStart])

	encryptionBegin := buffer.Len()

	buffer.AppendBytes(this.Version)
	buffer.Append(this.RequestIV)
	buffer.Append(this.RequestKey)
	buffer.AppendBytes(this.ResponseHeader, this.Option, byte(0), byte(0))
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

	timestampHash := md5.New()
	timestampHash.Write(hashTimestamp(timestamp))
	iv := timestampHash.Sum(nil)
	aesStream, err := v2crypto.NewAesEncryptionStream(this.User.ID.CmdKey(), iv)
	if err != nil {
		return nil, err
	}
	aesStream.XORKeyStream(buffer.Value[encryptionBegin:encryptionEnd], buffer.Value[encryptionBegin:encryptionEnd])

	return buffer, nil
}
