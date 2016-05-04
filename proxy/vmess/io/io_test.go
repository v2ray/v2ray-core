package io_test

import (
	"bytes"
	"crypto/rand"
	"io"
	"testing"

	"github.com/v2ray/v2ray-core/common/alloc"
	v2io "github.com/v2ray/v2ray-core/common/io"
	. "github.com/v2ray/v2ray-core/proxy/vmess/io"
	v2testing "github.com/v2ray/v2ray-core/testing"
	"github.com/v2ray/v2ray-core/testing/assert"
)

func TestAuthenticate(t *testing.T) {
	v2testing.Current(t)

	buffer := alloc.NewBuffer().Clear()
	buffer.AppendBytes(1, 2, 3, 4)
	Authenticate(buffer)
	assert.Bytes(buffer.Value).Equals([]byte{0, 8, 87, 52, 168, 125, 1, 2, 3, 4})

	b2, err := NewAuthChunkReader(buffer).Read()
	assert.Error(err).IsNil()
	assert.Bytes(b2.Value).Equals([]byte{1, 2, 3, 4})
}

func TestSingleIO(t *testing.T) {
	v2testing.Current(t)

	content := bytes.NewBuffer(make([]byte, 0, 1024*1024))

	writer := NewAuthChunkWriter(v2io.NewAdaptiveWriter(content))
	writer.Write(alloc.NewBuffer().Clear().AppendString("abcd"))
	writer.Release()

	reader := NewAuthChunkReader(content)
	buffer, err := reader.Read()
	assert.Error(err).IsNil()
	assert.Bytes(buffer.Value).Equals([]byte("abcd"))
}

func TestLargeIO(t *testing.T) {
	v2testing.Current(t)

	content := make([]byte, 1024*1024)
	rand.Read(content)

	chunckContent := bytes.NewBuffer(make([]byte, 0, len(content)*2))
	writer := NewAuthChunkWriter(v2io.NewAdaptiveWriter(chunckContent))
	writeSize := 0
	for {
		chunkSize := 7 * 1024
		if chunkSize+writeSize > len(content) {
			chunkSize = len(content) - writeSize
		}
		writer.Write(alloc.NewBuffer().Clear().Append(content[writeSize : writeSize+chunkSize]))
		writeSize += chunkSize
		if writeSize == len(content) {
			break
		}

		chunkSize = 8 * 1024
		if chunkSize+writeSize > len(content) {
			chunkSize = len(content) - writeSize
		}
		writer.Write(alloc.NewLargeBuffer().Clear().Append(content[writeSize : writeSize+chunkSize]))
		writeSize += chunkSize
		if writeSize == len(content) {
			break
		}

		chunkSize = 63 * 1024
		if chunkSize+writeSize > len(content) {
			chunkSize = len(content) - writeSize
		}
		writer.Write(alloc.NewLargeBuffer().Clear().Append(content[writeSize : writeSize+chunkSize]))
		writeSize += chunkSize
		if writeSize == len(content) {
			break
		}

		chunkSize = 64*1024 - 16
		if chunkSize+writeSize > len(content) {
			chunkSize = len(content) - writeSize
		}
		writer.Write(alloc.NewLargeBuffer().Clear().Append(content[writeSize : writeSize+chunkSize]))
		writeSize += chunkSize
		if writeSize == len(content) {
			break
		}
	}
	writer.Release()

	actualContent := make([]byte, 0, len(content))
	reader := NewAuthChunkReader(chunckContent)
	for {
		buffer, err := reader.Read()
		if err == io.EOF {
			break
		}
		assert.Error(err).IsNil()
		actualContent = append(actualContent, buffer.Value...)
	}

	assert.Int(len(actualContent)).Equals(len(content))
	assert.Bytes(actualContent).Equals(content)
}
