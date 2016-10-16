package encoding

import (
	"crypto/md5"
	"hash/fnv"
	"io"

	"v2ray.com/core/common/crypto"
	"v2ray.com/core/common/log"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/common/serial"
	"v2ray.com/core/proxy/vmess"
	"v2ray.com/core/transport"
)

type ServerSession struct {
	userValidator   protocol.UserValidator
	requestBodyKey  []byte
	requestBodyIV   []byte
	responseBodyKey []byte
	responseBodyIV  []byte
	responseHeader  byte
	responseWriter  io.Writer
}

// NewServerSession creates a new ServerSession, using the given UserValidator.
// The ServerSession instance doesn't take ownership of the validator.
func NewServerSession(validator protocol.UserValidator) *ServerSession {
	return &ServerSession{
		userValidator: validator,
	}
}

// Release implements common.Releaseable.
func (this *ServerSession) Release() {
	this.userValidator = nil
	this.requestBodyIV = nil
	this.requestBodyKey = nil
	this.responseBodyIV = nil
	this.responseBodyKey = nil
	this.responseWriter = nil
}

func (this *ServerSession) DecodeRequestHeader(reader io.Reader) (*protocol.RequestHeader, error) {
	buffer := make([]byte, 512)

	_, err := io.ReadFull(reader, buffer[:protocol.IDBytesLen])
	if err != nil {
		log.Info("Raw: Failed to read request header: ", err)
		return nil, io.EOF
	}

	user, timestamp, valid := this.userValidator.Get(buffer[:protocol.IDBytesLen])
	if !valid {
		return nil, protocol.ErrInvalidUser
	}

	timestampHash := md5.New()
	timestampHash.Write(hashTimestamp(timestamp))
	iv := timestampHash.Sum(nil)
	account, err := user.GetTypedAccount()
	if err != nil {
		log.Error("Vmess: Failed to get user account: ", err)
		return nil, err
	}

	aesStream := crypto.NewAesDecryptionStream(account.(*vmess.InternalAccount).ID.CmdKey(), iv)
	decryptor := crypto.NewCryptionReader(aesStream, reader)

	nBytes, err := io.ReadFull(decryptor, buffer[:41])
	if err != nil {
		log.Debug("Raw: Failed to read request header (", nBytes, " bytes): ", err)
		return nil, err
	}
	bufferLen := nBytes

	request := &protocol.RequestHeader{
		User:    user,
		Version: buffer[0],
	}

	if request.Version != Version {
		log.Info("Raw: Invalid protocol version ", request.Version)
		return nil, protocol.ErrInvalidVersion
	}

	this.requestBodyIV = append([]byte(nil), buffer[1:17]...)   // 16 bytes
	this.requestBodyKey = append([]byte(nil), buffer[17:33]...) // 16 bytes
	this.responseHeader = buffer[33]                            // 1 byte
	request.Option = protocol.RequestOption(buffer[34])         // 1 byte + 2 bytes reserved
	request.Command = protocol.RequestCommand(buffer[37])

	request.Port = v2net.PortFromBytes(buffer[38:40])

	switch buffer[40] {
	case AddrTypeIPv4:
		nBytes, err = io.ReadFull(decryptor, buffer[41:45]) // 4 bytes
		bufferLen += 4
		if err != nil {
			log.Debug("VMess: Failed to read target IPv4 (", nBytes, " bytes): ", err)
			return nil, err
		}
		request.Address = v2net.IPAddress(buffer[41:45])
	case AddrTypeIPv6:
		nBytes, err = io.ReadFull(decryptor, buffer[41:57]) // 16 bytes
		bufferLen += 16
		if err != nil {
			log.Debug("VMess: Failed to read target IPv6 (", nBytes, " bytes): ", nBytes, err)
			return nil, err
		}
		request.Address = v2net.IPAddress(buffer[41:57])
	case AddrTypeDomain:
		nBytes, err = io.ReadFull(decryptor, buffer[41:42])
		if err != nil {
			log.Debug("VMess: Failed to read target domain (", nBytes, " bytes): ", nBytes, err)
			return nil, err
		}
		domainLength := int(buffer[41])
		if domainLength == 0 {
			return nil, transport.ErrCorruptedPacket
		}
		nBytes, err = io.ReadFull(decryptor, buffer[42:42+domainLength])
		if err != nil {
			log.Debug("VMess: Failed to read target domain (", nBytes, " bytes): ", nBytes, err)
			return nil, err
		}
		bufferLen += 1 + domainLength
		request.Address = v2net.DomainAddress(string(buffer[42 : 42+domainLength]))
	}

	nBytes, err = io.ReadFull(decryptor, buffer[bufferLen:bufferLen+4])
	if err != nil {
		log.Debug("VMess: Failed to read checksum (", nBytes, " bytes): ", nBytes, err)
		return nil, err
	}

	fnv1a := fnv.New32a()
	fnv1a.Write(buffer[:bufferLen])
	actualHash := fnv1a.Sum32()
	expectedHash := serial.BytesToUint32(buffer[bufferLen : bufferLen+4])

	if actualHash != expectedHash {
		return nil, transport.ErrCorruptedPacket
	}

	return request, nil
}

func (this *ServerSession) DecodeRequestBody(reader io.Reader) io.Reader {
	aesStream := crypto.NewAesDecryptionStream(this.requestBodyKey, this.requestBodyIV)
	return crypto.NewCryptionReader(aesStream, reader)
}

func (this *ServerSession) EncodeResponseHeader(header *protocol.ResponseHeader, writer io.Writer) {
	responseBodyKey := md5.Sum(this.requestBodyKey)
	responseBodyIV := md5.Sum(this.requestBodyIV)
	this.responseBodyKey = responseBodyKey[:]
	this.responseBodyIV = responseBodyIV[:]

	aesStream := crypto.NewAesEncryptionStream(this.responseBodyKey, this.responseBodyIV)
	encryptionWriter := crypto.NewCryptionWriter(aesStream, writer)
	this.responseWriter = encryptionWriter

	encryptionWriter.Write([]byte{this.responseHeader, byte(header.Option)})
	err := MarshalCommand(header.Command, encryptionWriter)
	if err != nil {
		encryptionWriter.Write([]byte{0x00, 0x00})
	}
}

func (this *ServerSession) EncodeResponseBody(writer io.Writer) io.Writer {
	return this.responseWriter
}
