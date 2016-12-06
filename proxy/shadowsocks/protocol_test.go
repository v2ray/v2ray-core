package shadowsocks_test

import (
	"testing"

	"v2ray.com/core/common/alloc"
	"v2ray.com/core/common/loader"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/common/serial"
	. "v2ray.com/core/proxy/shadowsocks"
	"v2ray.com/core/testing/assert"
)

func TestUDPEncoding(t *testing.T) {
	assert := assert.On(t)

	request := &protocol.RequestHeader{
		Version: Version,
		Command: protocol.RequestCommandUDP,
		Address: v2net.LocalHostIP,
		Port:    1234,
		User: &protocol.User{
			Email: "love@v2ray.com",
			Account: loader.NewTypedSettings(&Account{
				Password:   "shadowsocks-password",
				CipherType: CipherType_AES_128_CFB,
				Ota:        Account_Disabled,
			}),
		},
	}

	data := alloc.NewLocalBuffer(256)
	data.AppendFunc(serial.WriteString("test string"))
	encodedData, err := EncodeUDPPacket(request, data)
	assert.Error(err).IsNil()

	decodedRequest, decodedData, err := DecodeUDPPacket(request.User, encodedData)
	assert.Error(err).IsNil()
	assert.Bytes(decodedData.Bytes()).Equals(data.Bytes())
	assert.Address(decodedRequest.Address).Equals(request.Address)
	assert.Port(decodedRequest.Port).Equals(request.Port)
}

func TestTCPRequest(t *testing.T) {
	assert := assert.On(t)

	request := &protocol.RequestHeader{
		Version: Version,
		Command: protocol.RequestCommandTCP,
		Address: v2net.LocalHostIP,
		Option:  RequestOptionOneTimeAuth,
		Port:    1234,
		User: &protocol.User{
			Email: "love@v2ray.com",
			Account: loader.NewTypedSettings(&Account{
				Password:   "tcp-password",
				CipherType: CipherType_CHACHA20,
			}),
		},
	}

	data := alloc.NewLocalBuffer(256)
	data.AppendFunc(serial.WriteString("test string"))
	cache := alloc.NewBuffer()

	writer, err := WriteTCPRequest(request, cache)
	assert.Error(err).IsNil()

	writer.Write(data)

	decodedRequest, reader, err := ReadTCPSession(request.User, cache)
	assert.Error(err).IsNil()
	assert.Address(decodedRequest.Address).Equals(request.Address)
	assert.Port(decodedRequest.Port).Equals(request.Port)

	decodedData, err := reader.Read()
	assert.Error(err).IsNil()
	assert.String(decodedData.String()).Equals("test string")
}

func TestUDPReaderWriter(t *testing.T) {
	assert := assert.On(t)

	user := &protocol.User{
		Account: loader.NewTypedSettings(&Account{
			Password:   "test-password",
			CipherType: CipherType_CHACHA20_IEFT,
		}),
	}
	cache := alloc.NewBuffer()
	writer := &UDPWriter{
		Writer: cache,
		Request: &protocol.RequestHeader{
			Version: Version,
			Address: v2net.DomainAddress("v2ray.com"),
			Port:    123,
			User:    user,
			Option:  RequestOptionOneTimeAuth,
		},
	}

	reader := &UDPReader{
		Reader: cache,
		User:   user,
	}

	b := alloc.NewBuffer()
	b.AppendFunc(serial.WriteString("test payload"))
	err := writer.Write(b)
	assert.Error(err).IsNil()

	payload, err := reader.Read()
	assert.Error(err).IsNil()
	assert.String(payload.String()).Equals("test payload")

	b = alloc.NewBuffer()
	b.AppendFunc(serial.WriteString("test payload 2"))
	err = writer.Write(b)
	assert.Error(err).IsNil()

	payload, err = reader.Read()
	assert.Error(err).IsNil()
	assert.String(payload.String()).Equals("test payload 2")
}
