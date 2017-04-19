package mux_test

import (
	"context"
	"testing"

	. "v2ray.com/core/app/proxyman/mux"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/net"
	"v2ray.com/core/testing/assert"
	"v2ray.com/core/transport/ray"
)

func TestReaderWriter(t *testing.T) {
	assert := assert.On(t)

	stream := ray.NewStream(context.Background())

	dest := net.TCPDestination(net.DomainAddress("v2ray.com"), 80)
	writer := NewWriter(1, dest, stream)

	dest2 := net.TCPDestination(net.LocalHostIP, 443)
	writer2 := NewWriter(2, dest2, stream)

	dest3 := net.TCPDestination(net.LocalHostIPv6, 18374)
	writer3 := NewWriter(3, dest3, stream)

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

	reader := NewReader(stream)
	meta, err := reader.ReadMetadata()
	assert.Error(err).IsNil()
	assert.Uint16(meta.SessionID).Equals(1)
	assert.Byte(byte(meta.SessionStatus)).Equals(byte(SessionStatusNew))
	assert.Destination(meta.Target).Equals(dest)
	assert.Byte(byte(meta.Option)).Equals(byte(OptionData))

	data, err := reader.Read()
	assert.Error(err).IsNil()
	assert.Int(len(data)).Equals(1)
	assert.String(data[0].String()).Equals("abcd")

	meta, err = reader.ReadMetadata()
	assert.Error(err).IsNil()
	assert.Byte(byte(meta.SessionStatus)).Equals(byte(SessionStatusNew))
	assert.Uint16(meta.SessionID).Equals(2)
	assert.Byte(byte(meta.Option)).Equals(0)
	assert.Destination(meta.Target).Equals(dest2)

	meta, err = reader.ReadMetadata()
	assert.Error(err).IsNil()
	assert.Byte(byte(meta.SessionStatus)).Equals(byte(SessionStatusKeep))
	assert.Uint16(meta.SessionID).Equals(1)
	assert.Byte(byte(meta.Option)).Equals(1)

	data, err = reader.Read()
	assert.Error(err).IsNil()
	assert.Int(len(data)).Equals(1)
	assert.String(data[0].String()).Equals("efgh")

	meta, err = reader.ReadMetadata()
	assert.Error(err).IsNil()
	assert.Byte(byte(meta.SessionStatus)).Equals(byte(SessionStatusNew))
	assert.Uint16(meta.SessionID).Equals(3)
	assert.Byte(byte(meta.Option)).Equals(1)
	assert.Destination(meta.Target).Equals(dest3)

	data, err = reader.Read()
	assert.Error(err).IsNil()
	assert.Int(len(data)).Equals(1)
	assert.String(data[0].String()).Equals("x")

	meta, err = reader.ReadMetadata()
	assert.Error(err).IsNil()
	assert.Byte(byte(meta.SessionStatus)).Equals(byte(SessionStatusEnd))
	assert.Uint16(meta.SessionID).Equals(1)
	assert.Byte(byte(meta.Option)).Equals(0)

	meta, err = reader.ReadMetadata()
	assert.Error(err).IsNil()
	assert.Byte(byte(meta.SessionStatus)).Equals(byte(SessionStatusEnd))
	assert.Uint16(meta.SessionID).Equals(3)
	assert.Byte(byte(meta.Option)).Equals(0)

	meta, err = reader.ReadMetadata()
	assert.Error(err).IsNil()
	assert.Byte(byte(meta.SessionStatus)).Equals(byte(SessionStatusKeep))
	assert.Uint16(meta.SessionID).Equals(2)
	assert.Byte(byte(meta.Option)).Equals(1)

	data, err = reader.Read()
	assert.Error(err).IsNil()
	assert.Int(len(data)).Equals(1)
	assert.String(data[0].String()).Equals("y")

	meta, err = reader.ReadMetadata()
	assert.Error(err).IsNil()
	assert.Byte(byte(meta.SessionStatus)).Equals(byte(SessionStatusEnd))
	assert.Uint16(meta.SessionID).Equals(2)
	assert.Byte(byte(meta.Option)).Equals(0)

	stream.Close()

	meta, err = reader.ReadMetadata()
	assert.Error(err).IsNotNil()
}
