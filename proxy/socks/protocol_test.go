package socks_test

import (
	"bytes"
	"testing"

	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/net"
	_ "v2ray.com/core/common/net/testing"
	"v2ray.com/core/common/protocol"
	. "v2ray.com/core/proxy/socks"
	. "v2ray.com/ext/assert"
)

func TestUDPEncoding(t *testing.T) {
	assert := With(t)

	b := buf.New()

	request := &protocol.RequestHeader{
		Address: net.IPAddress([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6}),
		Port:    1024,
	}
	writer := buf.NewSequentialWriter(NewUDPWriter(request, b))

	content := []byte{'a'}
	payload := buf.New()
	payload.Append(content)
	assert(writer.WriteMultiBuffer(buf.NewMultiBufferValue(payload)), IsNil)

	reader := NewUDPReader(b)

	decodedPayload, err := reader.ReadMultiBuffer()
	assert(err, IsNil)
	assert(decodedPayload[0].Bytes(), Equals, content)
}

func TestReadAddress(t *testing.T) {
	assert := With(t)

	data := []struct {
		AddrType byte
		Input    []byte
		Address  net.Address
		Port     net.Port
		Error    bool
	}{
		{
			AddrType: 0,
			Input:    []byte{0, 0, 0, 0},
			Error:    true,
		},
		{
			AddrType: 1,
			Input:    []byte{0, 0, 0, 0, 0, 53},
			Address:  net.IPAddress([]byte{0, 0, 0, 0}),
			Port:     net.Port(53),
		},
		{
			AddrType: 4,
			Input:    []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 0, 80},
			Address:  net.IPAddress([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6}),
			Port:     net.Port(80),
		},
		{
			AddrType: 3,
			Input:    []byte{9, 118, 50, 114, 97, 121, 46, 99, 111, 109, 0, 80},
			Address:  net.DomainAddress("v2ray.com"),
			Port:     net.Port(80),
		},
		{
			AddrType: 3,
			Input:    []byte{9, 118, 50, 114, 97, 121, 46, 99, 111, 109, 0},
			Error:    true,
		},
	}

	for _, tc := range data {
		b := buf.New()
		addr, port, err := ReadAddress(b, tc.AddrType, bytes.NewBuffer(tc.Input))
		b.Release()
		if tc.Error {
			assert(err, IsNotNil)
		} else {
			assert(addr, Equals, tc.Address)
			assert(port, Equals, tc.Port)
		}
	}
}
