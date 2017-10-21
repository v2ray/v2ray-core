package mux_test

import (
	"context"
	"io"
	"testing"

	. "v2ray.com/core/app/proxyman/mux"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/testing/assert"
	"v2ray.com/core/transport/ray"
)

func readAll(reader buf.Reader) (buf.MultiBuffer, error) {
	mb := buf.NewMultiBuffer()
	for {
		b, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		mb.AppendMulti(b)
	}
	return mb, nil
}

func TestReaderWriter(t *testing.T) {
	assert := assert.On(t)

	stream := ray.NewStream(context.Background())

	dest := net.TCPDestination(net.DomainAddress("v2ray.com"), 80)
	writer := NewWriter(1, dest, stream, protocol.TransferTypeStream)

	dest2 := net.TCPDestination(net.LocalHostIP, 443)
	writer2 := NewWriter(2, dest2, stream, protocol.TransferTypeStream)

	dest3 := net.TCPDestination(net.LocalHostIPv6, 18374)
	writer3 := NewWriter(3, dest3, stream, protocol.TransferTypeStream)

	writePayload := func(writer *Writer, payload ...byte) error {
		b := buf.New()
		b.Append(payload)
		return writer.Write(buf.NewMultiBufferValue(b))
	}

	assert.Error(writePayload(writer, 'a', 'b', 'c', 'd')).IsNil()
	assert.Error(writePayload(writer2)).IsNil()

	assert.Error(writePayload(writer, 'e', 'f', 'g', 'h')).IsNil()
	assert.Error(writePayload(writer3, 'x')).IsNil()

	writer.Close()
	writer3.Close()

	assert.Error(writePayload(writer2, 'y')).IsNil()
	writer2.Close()

	bytesReader := buf.ToBytesReader(stream)
	metaReader := NewMetadataReader(bytesReader)
	streamReader := NewStreamReader(bytesReader)

	meta, err := metaReader.Read()
	assert.Error(err).IsNil()
	assert.Uint16(meta.SessionID).Equals(1)
	assert.Byte(byte(meta.SessionStatus)).Equals(byte(SessionStatusNew))
	assert.Destination(meta.Target).Equals(dest)
	assert.Byte(byte(meta.Option)).Equals(byte(OptionData))

	data, err := readAll(streamReader)
	assert.Error(err).IsNil()
	assert.Int(len(data)).Equals(1)
	assert.String(data[0].String()).Equals("abcd")

	meta, err = metaReader.Read()
	assert.Error(err).IsNil()
	assert.Byte(byte(meta.SessionStatus)).Equals(byte(SessionStatusNew))
	assert.Uint16(meta.SessionID).Equals(2)
	assert.Byte(byte(meta.Option)).Equals(0)
	assert.Destination(meta.Target).Equals(dest2)

	meta, err = metaReader.Read()
	assert.Error(err).IsNil()
	assert.Byte(byte(meta.SessionStatus)).Equals(byte(SessionStatusKeep))
	assert.Uint16(meta.SessionID).Equals(1)
	assert.Byte(byte(meta.Option)).Equals(1)

	data, err = readAll(streamReader)
	assert.Error(err).IsNil()
	assert.Int(len(data)).Equals(1)
	assert.String(data[0].String()).Equals("efgh")

	meta, err = metaReader.Read()
	assert.Error(err).IsNil()
	assert.Byte(byte(meta.SessionStatus)).Equals(byte(SessionStatusNew))
	assert.Uint16(meta.SessionID).Equals(3)
	assert.Byte(byte(meta.Option)).Equals(1)
	assert.Destination(meta.Target).Equals(dest3)

	data, err = readAll(streamReader)
	assert.Error(err).IsNil()
	assert.Int(len(data)).Equals(1)
	assert.String(data[0].String()).Equals("x")

	meta, err = metaReader.Read()
	assert.Error(err).IsNil()
	assert.Byte(byte(meta.SessionStatus)).Equals(byte(SessionStatusEnd))
	assert.Uint16(meta.SessionID).Equals(1)
	assert.Byte(byte(meta.Option)).Equals(0)

	meta, err = metaReader.Read()
	assert.Error(err).IsNil()
	assert.Byte(byte(meta.SessionStatus)).Equals(byte(SessionStatusEnd))
	assert.Uint16(meta.SessionID).Equals(3)
	assert.Byte(byte(meta.Option)).Equals(0)

	meta, err = metaReader.Read()
	assert.Error(err).IsNil()
	assert.Byte(byte(meta.SessionStatus)).Equals(byte(SessionStatusKeep))
	assert.Uint16(meta.SessionID).Equals(2)
	assert.Byte(byte(meta.Option)).Equals(1)

	data, err = readAll(streamReader)
	assert.Error(err).IsNil()
	assert.Int(len(data)).Equals(1)
	assert.String(data[0].String()).Equals("y")

	meta, err = metaReader.Read()
	assert.Error(err).IsNil()
	assert.Byte(byte(meta.SessionStatus)).Equals(byte(SessionStatusEnd))
	assert.Uint16(meta.SessionID).Equals(2)
	assert.Byte(byte(meta.Option)).Equals(0)

	stream.Close()

	meta, err = metaReader.Read()
	assert.Error(err).IsNotNil()
	assert.Pointer(meta).IsNil()
}
