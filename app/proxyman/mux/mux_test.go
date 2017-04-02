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

	dest := net.TCPDestination(net.DomainAddress("v2ray.com"), 80)
	stream := ray.NewStream(context.Background())
	writer := NewWriter(1, dest, stream)

	payload := buf.New()
	payload.AppendBytes('a', 'b', 'c', 'd')
	assert.Error(writer.Write(payload)).IsNil()

	writer.Close()

	reader := NewReader(stream)
	meta, err := reader.ReadMetadata()
	assert.Error(err).IsNil()
	assert.Uint16(meta.SessionID).Equals(1)
	assert.Byte(byte(meta.SessionStatus)).Equals(byte(SessionStatusNew))
	assert.Destination(meta.Target).Equals(dest)
	assert.Byte(byte(meta.Option)).Equals(byte(OptionData))

	data, more, err := reader.Read()
	assert.Error(err).IsNil()
	assert.Bool(more).IsFalse()
	assert.String(data.String()).Equals("abcd")

	meta, err = reader.ReadMetadata()
	assert.Error(err).IsNil()
	assert.Byte(byte(meta.SessionStatus)).Equals(byte(SessionStatusEnd))
	assert.Uint16(meta.SessionID).Equals(1)
}
