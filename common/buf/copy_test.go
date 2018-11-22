package buf_test

import (
	"crypto/rand"
	"io"
	"testing"

	"github.com/golang/mock/gomock"

	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/errors"
	"v2ray.com/core/testing/mocks"
)

func TestReadError(t *testing.T) {
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	mockReader := mocks.NewReader(mockCtl)
	mockReader.EXPECT().Read(gomock.Any()).Return(0, errors.New("error"))

	err := buf.Copy(buf.NewReader(mockReader), buf.Discard)
	if err == nil {
		t.Fatal("expected error, but nil")
	}

	if !buf.IsReadError(err) {
		t.Error("expected to be ReadError, but not")
	}

	if err.Error() != "error" {
		t.Fatal("unexpected error message: ", err.Error())
	}
}

func TestWriteError(t *testing.T) {
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	mockWriter := mocks.NewWriter(mockCtl)
	mockWriter.EXPECT().Write(gomock.Any()).Return(0, errors.New("error"))

	err := buf.Copy(buf.NewReader(rand.Reader), buf.NewWriter(mockWriter))
	if err == nil {
		t.Fatal("expected error, but nil")
	}

	if !buf.IsWriteError(err) {
		t.Error("expected to be WriteError, but not")
	}

	if err.Error() != "error" {
		t.Fatal("unexpected error message: ", err.Error())
	}
}

type TestReader struct{}

func (TestReader) Read(b []byte) (int, error) {
	return len(b), nil
}

func BenchmarkCopy(b *testing.B) {
	reader := buf.NewReader(io.LimitReader(TestReader{}, 10240))
	writer := buf.Discard

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = buf.Copy(reader, writer)
	}
}
