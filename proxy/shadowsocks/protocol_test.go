package shadowsocks_test

import (
	"testing"

	"v2ray.com/core/common/buf"
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
			Account: serial.ToTypedMessage(&Account{
				Password:   "shadowsocks-password",
				CipherType: CipherType_AES_128_CFB,
				Ota:        Account_Disabled,
			}),
		},
	}

	data := buf.NewLocal(256)
	data.AppendSupplier(serial.WriteString("test string"))
	encodedData, err := EncodeUDPPacket(request, data.Bytes())
	assert.Error(err).IsNil()

	decodedRequest, decodedData, err := DecodeUDPPacket(request.User, encodedData)
	assert.Error(err).IsNil()
	assert.Bytes(decodedData.Bytes()).Equals(data.Bytes())
	assert.Address(decodedRequest.Address).Equals(request.Address)
	assert.Port(decodedRequest.Port).Equals(request.Port)
}

func TestTCPRequest(t *testing.T) {
	assert := assert.On(t)

	cases := []struct {
		request *protocol.RequestHeader
		payload []byte
	}{
		{
			request: &protocol.RequestHeader{
				Version: Version,
				Command: protocol.RequestCommandTCP,
				Address: v2net.LocalHostIP,
				Option:  RequestOptionOneTimeAuth,
				Port:    1234,
				User: &protocol.User{
					Email: "love@v2ray.com",
					Account: serial.ToTypedMessage(&Account{
						Password:   "tcp-password",
						CipherType: CipherType_CHACHA20,
					}),
				},
			},
			payload: []byte("test string"),
		},
		{
			request: &protocol.RequestHeader{
				Version: Version,
				Command: protocol.RequestCommandTCP,
				Address: v2net.LocalHostIPv6,
				Option:  RequestOptionOneTimeAuth,
				Port:    1234,
				User: &protocol.User{
					Email: "love@v2ray.com",
					Account: serial.ToTypedMessage(&Account{
						Password:   "password",
						CipherType: CipherType_AES_256_CFB,
					}),
				},
			},
			payload: []byte("test string"),
		},
		{
			request: &protocol.RequestHeader{
				Version: Version,
				Command: protocol.RequestCommandTCP,
				Address: v2net.DomainAddress("v2ray.com"),
				Option:  RequestOptionOneTimeAuth,
				Port:    1234,
				User: &protocol.User{
					Email: "love@v2ray.com",
					Account: serial.ToTypedMessage(&Account{
						Password:   "password",
						CipherType: CipherType_CHACHA20_IETF,
					}),
				},
			},
			payload: []byte("test string"),
		},
	}

	runTest := func(request *protocol.RequestHeader, payload []byte) {
		data := buf.New()
		defer data.Release()
		data.Append(payload)

		cache := buf.New()
		defer cache.Release()

		writer, err := WriteTCPRequest(request, cache)
		assert.Error(err).IsNil()

		assert.Error(writer.Write(buf.NewMultiBufferValue(data))).IsNil()

		decodedRequest, reader, err := ReadTCPSession(request.User, cache)
		assert.Error(err).IsNil()
		assert.Address(decodedRequest.Address).Equals(request.Address)
		assert.Port(decodedRequest.Port).Equals(request.Port)

		decodedData, err := reader.Read()
		assert.Error(err).IsNil()
		assert.String(decodedData[0].String()).Equals(string(payload))
	}

	for _, test := range cases {
		runTest(test.request, test.payload)
	}

}

func TestUDPReaderWriter(t *testing.T) {
	assert := assert.On(t)

	user := &protocol.User{
		Account: serial.ToTypedMessage(&Account{
			Password:   "test-password",
			CipherType: CipherType_CHACHA20_IETF,
		}),
	}
	cache := buf.New()
	writer := buf.NewSequentialWriter(&UDPWriter{
		Writer: cache,
		Request: &protocol.RequestHeader{
			Version: Version,
			Address: v2net.DomainAddress("v2ray.com"),
			Port:    123,
			User:    user,
			Option:  RequestOptionOneTimeAuth,
		},
	})

	reader := &UDPReader{
		Reader: cache,
		User:   user,
	}

	b := buf.New()
	b.AppendSupplier(serial.WriteString("test payload"))
	err := writer.Write(buf.NewMultiBufferValue(b))
	assert.Error(err).IsNil()

	payload, err := reader.Read()
	assert.Error(err).IsNil()
	assert.String(payload[0].String()).Equals("test payload")

	b = buf.New()
	b.AppendSupplier(serial.WriteString("test payload 2"))
	err = writer.Write(buf.NewMultiBufferValue(b))
	assert.Error(err).IsNil()

	payload, err = reader.Read()
	assert.Error(err).IsNil()
	assert.String(payload[0].String()).Equals("test payload 2")
}
