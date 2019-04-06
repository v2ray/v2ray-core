package shadowsocks_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/protocol"
	. "v2ray.com/core/proxy/shadowsocks"
)

func toAccount(a *Account) protocol.Account {
	account, err := a.AsAccount()
	common.Must(err)
	return account
}

func TestUDPEncoding(t *testing.T) {
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
	common.Must2(data.WriteString("test string"))
	encodedData, err := EncodeUDPPacket(request, data.Bytes())
	common.Must(err)

	decodedRequest, decodedData, err := DecodeUDPPacket(request.User, encodedData)
	common.Must(err)

	if r := cmp.Diff(decodedData.Bytes(), data.Bytes()); r != "" {
		t.Error("data: ", r)
	}

	if r := cmp.Diff(decodedRequest, request); r != "" {
		t.Error("request: ", r)
	}
}

func TestTCPRequest(t *testing.T) {
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
		common.Must(err)

		common.Must(writer.WriteMultiBuffer(buf.MultiBuffer{data}))

		decodedRequest, reader, err := ReadTCPSession(request.User, cache)
		common.Must(err)
		if r := cmp.Diff(decodedRequest, request); r != "" {
			t.Error("request: ", r)
		}

		decodedData, err := reader.ReadMultiBuffer()
		common.Must(err)
		if r := cmp.Diff(decodedData[0].Bytes(), payload); r != "" {
			t.Error("data: ", r)
		}
	}

	for _, test := range cases {
		runTest(test.request, test.payload)
	}

}

func TestUDPReaderWriter(t *testing.T) {
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
		common.Must2(b.WriteString("test payload"))
		common.Must(writer.WriteMultiBuffer(buf.MultiBuffer{b}))

		payload, err := reader.ReadMultiBuffer()
		common.Must(err)
		if payload[0].String() != "test payload" {
			t.Error("unexpected output: ", payload[0].String())
		}
	}

	{
		b := buf.New()
		common.Must2(b.WriteString("test payload 2"))
		common.Must(writer.WriteMultiBuffer(buf.MultiBuffer{b}))

		payload, err := reader.ReadMultiBuffer()
		common.Must(err)
		if payload[0].String() != "test payload 2" {
			t.Error("unexpected output: ", payload[0].String())
		}
	}
}
