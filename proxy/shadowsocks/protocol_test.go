package shadowsocks_test

import (
	"testing"

	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/common/serial"
	. "v2ray.com/core/proxy/shadowsocks"
	. "v2ray.com/ext/assert"
)

func toAccount(a *Account) protocol.Account {
	account, err := a.AsAccount()
	common.Must(err)
	return account
}

func TestUDPEncoding(t *testing.T) {
	assert := With(t)

	request := &protocol.RequestHeader{
		Version: Version,
		Command: protocol.RequestCommandUDP,
		Address: net.LocalHostIP,
		Port:    1234,
		User: &protocol.MemoryUser{
			Email: "love@v2ray.com",
			Account: toAccount(&Account{
				Password:   "shadowsocks-password",
				CipherType: CipherType_AES_128_CFB,
				Ota:        Account_Disabled,
			}),
		},
	}

	data := buf.New()
	data.AppendSupplier(serial.WriteString("test string"))
	encodedData, err := EncodeUDPPacket(request, data.Bytes())
	assert(err, IsNil)

	decodedRequest, decodedData, err := DecodeUDPPacket(request.User, encodedData)
	assert(err, IsNil)
	assert(decodedData.Bytes(), Equals, data.Bytes())
	assert(decodedRequest.Address, Equals, request.Address)
	assert(decodedRequest.Port, Equals, request.Port)
	assert(decodedRequest.Command, Equals, request.Command)
}

func TestTCPRequest(t *testing.T) {
	assert := With(t)

	cases := []struct {
		request *protocol.RequestHeader
		payload []byte
	}{
		{
			request: &protocol.RequestHeader{
				Version: Version,
				Command: protocol.RequestCommandTCP,
				Address: net.LocalHostIP,
				Option:  RequestOptionOneTimeAuth,
				Port:    1234,
				User: &protocol.MemoryUser{
					Email: "love@v2ray.com",
					Account: toAccount(&Account{
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
				Address: net.LocalHostIPv6,
				Option:  RequestOptionOneTimeAuth,
				Port:    1234,
				User: &protocol.MemoryUser{
					Email: "love@v2ray.com",
					Account: toAccount(&Account{
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
				Address: net.DomainAddress("v2ray.com"),
				Option:  RequestOptionOneTimeAuth,
				Port:    1234,
				User: &protocol.MemoryUser{
					Email: "love@v2ray.com",
					Account: toAccount(&Account{
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
		common.Must2(data.Write(payload))

		cache := buf.New()
		defer cache.Release()

		writer, err := WriteTCPRequest(request, cache)
		assert(err, IsNil)

		assert(writer.WriteMultiBuffer(buf.NewMultiBufferValue(data)), IsNil)

		decodedRequest, reader, err := ReadTCPSession(request.User, cache)
		assert(err, IsNil)
		assert(decodedRequest.Address, Equals, request.Address)
		assert(decodedRequest.Port, Equals, request.Port)
		assert(decodedRequest.Command, Equals, request.Command)

		decodedData, err := reader.ReadMultiBuffer()
		assert(err, IsNil)
		assert(decodedData[0].String(), Equals, string(payload))
	}

	for _, test := range cases {
		runTest(test.request, test.payload)
	}

}

func TestUDPReaderWriter(t *testing.T) {
	assert := With(t)

	user := &protocol.MemoryUser{
		Account: toAccount(&Account{
			Password:   "test-password",
			CipherType: CipherType_CHACHA20_IETF,
		}),
	}
	cache := buf.New()
	defer cache.Release()

	writer := &buf.SequentialWriter{Writer: &UDPWriter{
		Writer: cache,
		Request: &protocol.RequestHeader{
			Version: Version,
			Address: net.DomainAddress("v2ray.com"),
			Port:    123,
			User:    user,
			Option:  RequestOptionOneTimeAuth,
		},
	}}

	reader := &UDPReader{
		Reader: cache,
		User:   user,
	}

	{
		b := buf.New()
		b.AppendSupplier(serial.WriteString("test payload"))
		err := writer.WriteMultiBuffer(buf.NewMultiBufferValue(b))
		assert(err, IsNil)

		payload, err := reader.ReadMultiBuffer()
		assert(err, IsNil)
		assert(payload[0].String(), Equals, "test payload")
	}

	{
		b := buf.New()
		b.AppendSupplier(serial.WriteString("test payload 2"))
		err := writer.WriteMultiBuffer(buf.NewMultiBufferValue(b))
		assert(err, IsNil)

		payload, err := reader.ReadMultiBuffer()
		assert(err, IsNil)
		assert(payload[0].String(), Equals, "test payload 2")
	}
}
