package io

import (
	"errors"
	"hash"
	"hash/fnv"
	"io"
	"v2ray.com/core/common/alloc"
	"v2ray.com/core/common/serial"
)

// Private: Visible for testing.
type Validator struct {
	actualAuth   hash.Hash32
	expectedAuth uint32
}

func NewValidator(expectedAuth uint32) *Validator {
	return &Validator{
		actualAuth:   fnv.New32a(),
		expectedAuth: expectedAuth,
	}
}

func (v *Validator) Consume(b []byte) {
	v.actualAuth.Write(b)
}

func (v *Validator) Validate() bool {
	return v.actualAuth.Sum32() == v.expectedAuth
}

type AuthChunkReader struct {
	reader      io.Reader
	last        *alloc.Buffer
	chunkLength int
	validator   *Validator
}

func NewAuthChunkReader(reader io.Reader) *AuthChunkReader {
	return &AuthChunkReader{
		reader:      reader,
		chunkLength: -1,
	}
}

func (v *AuthChunkReader) Read() (*alloc.Buffer, error) {
	var buffer *alloc.Buffer
	if v.last != nil {
		buffer = v.last
		v.last = nil
	} else {
		buffer = alloc.NewBuffer().Clear()
	}

	if v.chunkLength == -1 {
		for buffer.Len() < 6 {
			_, err := buffer.FillFrom(v.reader)
			if err != nil {
				buffer.Release()
				return nil, io.ErrUnexpectedEOF
			}
		}
		length := serial.BytesToUint16(buffer.Value[:2])
		v.chunkLength = int(length) - 4
		v.validator = NewValidator(serial.BytesToUint32(buffer.Value[2:6]))
		buffer.SliceFrom(6)
		if buffer.Len() < v.chunkLength && v.chunkLength <= 2048 {
			_, err := buffer.FillFrom(v.reader)
			if err != nil {
				buffer.Release()
				return nil, io.ErrUnexpectedEOF
			}
		}
	} else if buffer.Len() < v.chunkLength {
		_, err := buffer.FillFrom(v.reader)
		if err != nil {
			buffer.Release()
			return nil, io.ErrUnexpectedEOF
		}
	}

	if v.chunkLength == 0 {
		buffer.Release()
		return nil, io.EOF
	}

	if buffer.Len() < v.chunkLength {
		v.validator.Consume(buffer.Value)
		v.chunkLength -= buffer.Len()
	} else {
		v.validator.Consume(buffer.Value[:v.chunkLength])
		if !v.validator.Validate() {
			buffer.Release()
			return nil, errors.New("VMess|AuthChunkReader: Invalid auth.")
		}
		leftLength := buffer.Len() - v.chunkLength
		if leftLength > 0 {
			v.last = alloc.NewBuffer().Clear()
			v.last.Append(buffer.Value[v.chunkLength:])
			buffer.Slice(0, v.chunkLength)
		}

		v.chunkLength = -1
		v.validator = nil
	}

	return buffer, nil
}

func (v *AuthChunkReader) Release() {
	v.reader = nil
	v.last.Release()
	v.last = nil
	v.validator = nil
}
