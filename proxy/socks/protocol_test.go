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
	writer := &buf.SequentialWriter{Writer: NewUDPWriter(request, b)}

	content := []byte{'a'}
	payload := buf.New()
	payload.Write(content)
	assert(writer.WriteMultiBuffer(buf.NewMultiBufferValue(payload)), IsNil)

	reader := NewUDPReader(b)

	decodedPayload, err := reader.ReadMultiBuffer()
	assert(err, IsNil)
	assert(decodedPayload[0].Bytes(), Equals, content)
}

func TestReadUsernamePassword(t *testing.T) {
	testCases := []struct {
		Input    []byte
		Username string
		Password string
		Error    bool
	}{
		{
			Input:    []byte{0x05, 0x01, 'a', 0x02, 'b', 'c'},
			Username: "a",
			Password: "bc",
		},
		{
			Input: []byte{0x05, 0x18, 'a', 0x02, 'b', 'c'},
			Error: true,
		},
	}

	for _, testCase := range testCases {
		reader := bytes.NewReader(testCase.Input)
		username, password, err := ReadUsernamePassword(reader)
		if testCase.Error {
			if err == nil {
				t.Error("for input: ", testCase.Input, " expect error, but actually nil")
			}
		} else {
			if err != nil {
				t.Error("for input: ", testCase.Input, " expect no error, but actually ", err.Error())
			}
			if testCase.Username != username {
				t.Error("for input: ", testCase.Input, " expect username ", testCase.Username, " but actually ", username)
			}
			if testCase.Password != password {
				t.Error("for input: ", testCase.Input, " expect passowrd ", testCase.Password, " but actually ", password)
			}
		}
	}
}

func TestReadUntilNull(t *testing.T) {
	testCases := []struct {
		Input  []byte
		Output string
		Error  bool
	}{
		{
			Input:  []byte{'a', 'b', 0x00},
			Output: "ab",
		},
		{
			Input: []byte{'a'},
			Error: true,
		},
	}

	for _, testCase := range testCases {
		reader := bytes.NewReader(testCase.Input)
		value, err := ReadUntilNull(reader)
		if testCase.Error {
			if err == nil {
				t.Error("for input: ", testCase.Input, " expect error, but actually nil")
			}
		} else {
			if err != nil {
				t.Error("for input: ", testCase.Input, " expect no error, but actually ", err.Error())
			}
			if testCase.Output != value {
				t.Error("for input: ", testCase.Input, " expect output ", testCase.Output, " but actually ", value)
			}
		}
	}
}
