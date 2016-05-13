package io

import (
	"hash"
	"hash/fnv"
	"io"

	"github.com/v2ray/v2ray-core/common/alloc"
	"github.com/v2ray/v2ray-core/common/log"
	"github.com/v2ray/v2ray-core/common/serial"
	"github.com/v2ray/v2ray-core/transport"
)

// @Private
func AllocBuffer(size int) *alloc.Buffer {
	if size < 8*1024-16 {
		return alloc.NewBuffer()
	}
	return alloc.NewLargeBuffer()
}

// @Private
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

func (this *Validator) Consume(b []byte) {
	this.actualAuth.Write(b)
}

func (this *Validator) Validate() bool {
	log.Debug("VMess Reader: Expected auth ", this.expectedAuth, " actual auth: ", this.actualAuth.Sum32())
	return this.actualAuth.Sum32() == this.expectedAuth
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

func (this *AuthChunkReader) Read() (*alloc.Buffer, error) {
	var buffer *alloc.Buffer
	if this.last != nil {
		buffer = this.last
		this.last = nil
	} else {
		buffer = AllocBuffer(this.chunkLength).Clear()
	}

	if this.chunkLength == -1 {
		for buffer.Len() < 6 {
			_, err := buffer.FillFrom(this.reader)
			if err != nil {
				buffer.Release()
				return nil, err
			}
		}
		log.Debug("VMess Reader: raw buffer: ", buffer.Value)
		length := serial.BytesLiteral(buffer.Value[:2]).Uint16Value()
		this.chunkLength = int(length) - 4
		this.validator = NewValidator(serial.BytesLiteral(buffer.Value[2:6]).Uint32Value())
		buffer.SliceFrom(6)
	} else if buffer.Len() < this.chunkLength {
		_, err := buffer.FillFrom(this.reader)
		if err != nil {
			buffer.Release()
			return nil, err
		}
	}

	if this.chunkLength == 0 {
		buffer.Release()
		return nil, io.EOF
	}

	if buffer.Len() < this.chunkLength {
		this.validator.Consume(buffer.Value)
		this.chunkLength -= buffer.Len()
	} else {
		this.validator.Consume(buffer.Value[:this.chunkLength])
		if !this.validator.Validate() {
			buffer.Release()
			return nil, transport.ErrorCorruptedPacket
		}
		leftLength := buffer.Len() - this.chunkLength
		if leftLength > 0 {
			this.last = AllocBuffer(leftLength).Clear()
			this.last.Append(buffer.Value[this.chunkLength:])
			buffer.Slice(0, this.chunkLength)
		}

		this.chunkLength = -1
		this.validator = nil
	}

	return buffer, nil
}

func (this *AuthChunkReader) Release() {
	this.reader = nil
	this.last.Release()
	this.last = nil
	this.validator = nil
}
