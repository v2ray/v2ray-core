package socks_test

import (
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"

	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/protocol"
	. "v2ray.com/core/proxy/socks"
)

func TestUDPEncoding(t *testing.T) {
	b := buf.New()

	request := &protocol.RequestHeader{
		Address: net.IPAddress([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6}),
		Port:    1024,
	}
	writer := &buf.SequentialWriter{Writer: NewUDPWriter(request, b)}

	content := []byte{'a'}
	payload := buf.New()
	payload.Write(content)
	common.Must(writer.WriteMultiBuffer(buf.MultiBuffer{payload}))

	reader := NewUDPReader(b)

	decodedPayload, err := reader.ReadMultiBuffer()
	common.Must(err)
	if r := cmp.Diff(decodedPayload[0].Bytes(), content); r != "" {
		t.Error(r)
	}
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

func BenchmarkReadUsernamePassword(b *testing.B) {
	input := []byte{0x05, 0x01, 'a', 0x02, 'b', 'c'}
	buffer := buf.New()
	buffer.Write(input)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := ReadUsernamePassword(buffer)
		common.Must(err)
		buffer.Clear()
		buffer.Extend(int32(len(input)))
	}
}
