package tls_test

import (
	"context"
	"testing"

	"v2ray.com/core/common/buf"
	. "v2ray.com/core/transport/internet/headers/tls"
	. "v2ray.com/ext/assert"
)

func TestDTLSWrite(t *testing.T) {
	assert := With(t)

	content := []byte{'a', 'b', 'c', 'd', 'e', 'f', 'g'}
	dtlsRaw, err := New(context.Background(), &PacketConfig{})
	assert(err, IsNil)

	dtls := dtlsRaw.(*DTLS)

	payload := buf.New()
	dtls.Serialize(payload.Extend(dtls.Size()))
	payload.Write(content)

	assert(payload.Len(), Equals, int32(len(content))+dtls.Size())
}
