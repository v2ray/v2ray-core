package encoding

import (
	"crypto/md5"
	"errors"
	"hash/fnv"
	"io"
	"v2ray.com/core/common/crypto"
	"v2ray.com/core/common/log"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/common/serial"
	"v2ray.com/core/proxy/vmess"
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
func (v *ServerSession) Release() {
	v.userValidator = nil
	v.requestBodyIV = nil
	v.requestBodyKey = nil
	v.responseBodyIV = nil
	v.responseBodyKey = nil
	v.responseWriter = nil
}

func (v *ServerSession) DecodeRequestHeader(reader io.Reader) (*protocol.RequestHeader, error) {
	buffer := make([]byte, 512)

	_, err := io.ReadFull(reader, buffer[:protocol.IDBytesLen])
	if err != nil {
		log.Info("Raw: Failed to read request header: ", err)
		return nil, io.EOF
	}

	user, timestamp, valid := v.userValidator.Get(buffer[:protocol.IDBytesLen])
	if !valid {
		return nil, protocol.ErrInvalidUser
	}

	timestampHash := md5.New()
	timestampHash.Write(hashTimestamp(timestamp))
	iv := timestampHash.Sum(nil)
	account, err := user.GetTypedAccount()
	if err != nil {
		return nil, errors.New("VMess|Server: Failed to get user account: " + err.Error())
	}

	aesStream := crypto.NewAesDecryptionStream(account.(*vmess.InternalAccount).ID.CmdKey(), iv)
	decryptor := crypto.NewCryptionReader(aesStream, reader)

	nBytes, err := io.ReadFull(decryptor, buffer[:41])
	if err != nil {
		return nil, errors.New("VMess|Server: Failed to read request header: " + err.Error())
	}
	bufferLen := nBytes

	request := &protocol.RequestHeader{
		User:    user,
		Version: buffer[0],
	}

	if request.Version != Version {
		log.Info("VMess|Server: Invalid protocol version ", request.Version)
		return nil, protocol.ErrInvalidVersion
	}

	v.requestBodyIV = append([]byte(nil), buffer[1:17]...)   // 16 bytes
	v.requestBodyKey = append([]byte(nil), buffer[17:33]...) // 16 bytes
	v.responseHeader = buffer[33]                            // 1 byte
	request.Option = protocol.RequestOption(buffer[34])         // 1 byte + 2 bytes reserved
	request.Command = protocol.RequestCommand(buffer[37])

	request.Port = v2net.PortFromBytes(buffer[38:40])

	switch buffer[40] {
	case AddrTypeIPv4:
		nBytes, err = io.ReadFull(decryptor, buffer[41:45]) // 4 bytes
		bufferLen += 4
		if err != nil {
			return nil, errors.New("VMess|Server: Failed to read IPv4: " + err.Error())
		}
		request.Address = v2net.IPAddress(buffer[41:45])
	case AddrTypeIPv6:
		nBytes, err = io.ReadFull(decryptor, buffer[41:57]) // 16 bytes
		bufferLen += 16
		if err != nil {
			return nil, errors.New("VMess|Server: Failed to read IPv6 address: " + err.Error())
		}
		request.Address = v2net.IPAddress(buffer[41:57])
	case AddrTypeDomain:
		_, err = io.ReadFull(decryptor, buffer[41:42])
		if err != nil {
			return nil, errors.New("VMess:Server: Failed to read domain: " + err.Error())
		}
		domainLength := int(buffer[41])
		if domainLength == 0 {
			return nil, errors.New("VMess|Server: Zero domain length.")
		}
		nBytes, err = io.ReadFull(decryptor, buffer[42:42+domainLength])
		if err != nil {
			return nil, errors.New("VMess|Server: Failed to read domain: " + err.Error())
		}
		bufferLen += 1 + domainLength
		request.Address = v2net.DomainAddress(string(buffer[42 : 42+domainLength]))
	}

	nBytes, err = io.ReadFull(decryptor, buffer[bufferLen:bufferLen+4])
	if err != nil {
		return nil, errors.New("VMess|Server: Failed to read checksum: " + err.Error())
	}

	fnv1a := fnv.New32a()
	fnv1a.Write(buffer[:bufferLen])
	actualHash := fnv1a.Sum32()
	expectedHash := serial.BytesToUint32(buffer[bufferLen : bufferLen+4])

	if actualHash != expectedHash {
		return nil, errors.New("VMess|Server: Invalid auth.")
	}

	return request, nil
}

func (v *ServerSession) DecodeRequestBody(reader io.Reader) io.Reader {
	aesStream := crypto.NewAesDecryptionStream(v.requestBodyKey, v.requestBodyIV)
	return crypto.NewCryptionReader(aesStream, reader)
}

func (v *ServerSession) EncodeResponseHeader(header *protocol.ResponseHeader, writer io.Writer) {
	responseBodyKey := md5.Sum(v.requestBodyKey)
	responseBodyIV := md5.Sum(v.requestBodyIV)
	v.responseBodyKey = responseBodyKey[:]
	v.responseBodyIV = responseBodyIV[:]

	aesStream := crypto.NewAesEncryptionStream(v.responseBodyKey, v.responseBodyIV)
	encryptionWriter := crypto.NewCryptionWriter(aesStream, writer)
	v.responseWriter = encryptionWriter

	encryptionWriter.Write([]byte{v.responseHeader, byte(header.Option)})
	err := MarshalCommand(header.Command, encryptionWriter)
	if err != nil {
		encryptionWriter.Write([]byte{0x00, 0x00})
	}
}

func (v *ServerSession) EncodeResponseBody(writer io.Writer) io.Writer {
	return v.responseWriter
}
