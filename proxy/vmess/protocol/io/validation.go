package io

import (
	"errors"
	"hash/fnv"
	"io"

	"github.com/v2ray/v2ray-core/common/alloc"
	"github.com/v2ray/v2ray-core/transport"
)

var (
	TruncatedPayload = errors.New("Truncated payload.")
)

type ValidationReader struct {
	reader io.Reader
	buffer *alloc.Buffer
}

func NewValidationReader(reader io.Reader) *ValidationReader {
	return &ValidationReader{
		reader: reader,
		buffer: alloc.NewLargeBuffer().Clear(),
	}
}

func (this *ValidationReader) Read(data []byte) (int, error) {
	nBytes, err := this.reader.Read(data)
	if err != nil {
		return nBytes, err
	}
	nBytesActual := 0
	dataActual := data[:]
	for {
		payload, rest, err := parsePayload(data)
		if err != nil {
			return nBytesActual, err
		}
		copy(dataActual, payload)
		nBytesActual += len(payload)
		dataActual = dataActual[nBytesActual:]
		if len(rest) == 0 {
			break
		}
		data = rest
	}
	return nBytesActual, nil
}

func parsePayload(data []byte) (payload []byte, rest []byte, err error) {
	dataLen := len(data)
	if dataLen < 6 {
		err = TruncatedPayload
		return
	}
	payloadLen := int(data[0])<<8 + int(data[1])
	if dataLen < payloadLen+6 {
		err = TruncatedPayload
		return
	}

	payload = data[6 : 6+payloadLen]
	rest = data[6+payloadLen:]

	fnv1a := fnv.New32a()
	fnv1a.Write(payload)
	actualHash := fnv1a.Sum32()
	expectedHash := uint32(data[2])<<24 + uint32(data[3])<<16 + uint32(data[4])<<8 + uint32(data[5])
	if actualHash != expectedHash {
		err = transport.CorruptedPacket
		return
	}
	return
}
