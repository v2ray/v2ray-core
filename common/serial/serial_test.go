package serial_test

import (
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"

	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/serial"
)

func TestUint16Serial(t *testing.T) {
	b := buf.New()
	defer b.Release()

	n, err := serial.WriteUint16(b, 10)
	common.Must(err)
	if n != 2 {
		t.Error("expect 2 bytes writtng, but actually ", n)
	}
	if diff := cmp.Diff(b.Bytes(), []byte{0, 10}); diff != "" {
		t.Error(diff)
	}
}

func TestUint64Serial(t *testing.T) {
	b := buf.New()
	defer b.Release()

	n, err := serial.WriteUint64(b, 10)
	common.Must(err)
	if n != 8 {
		t.Error("expect 8 bytes writtng, but actually ", n)
	}
	if diff := cmp.Diff(b.Bytes(), []byte{0, 0, 0, 0, 0, 0, 0, 10}); diff != "" {
		t.Error(diff)
	}
}

func TestReadUint16(t *testing.T) {
	testCases := []struct {
		Input  []byte
		Output uint16
	}{
		{
			Input:  []byte{0, 1},
			Output: 1,
		},
	}

	for _, testCase := range testCases {
		v, err := serial.ReadUint16(bytes.NewReader(testCase.Input))
		common.Must(err)
		if v != testCase.Output {
			t.Error("for input ", testCase.Input, " expect output ", testCase.Output, " but got ", v)
		}
	}
}

func BenchmarkReadUint16(b *testing.B) {
	reader := buf.New()
	defer reader.Release()

	common.Must2(reader.Write([]byte{0, 1}))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := serial.ReadUint16(reader)
		common.Must(err)
		reader.Clear()
		reader.Extend(2)
	}
}

func BenchmarkWriteUint64(b *testing.B) {
	writer := buf.New()
	defer writer.Release()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := serial.WriteUint64(writer, 8)
		common.Must(err)
		writer.Clear()
	}
}
