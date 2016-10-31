package http_test

import (
	"testing"

	"v2ray.com/core/common/alloc"
	"v2ray.com/core/testing/assert"
	. "v2ray.com/core/transport/internet/authenticators/http"
)

func TestRequestOpenSeal(t *testing.T) {
	assert := assert.On(t)

	content := []byte{'a', 'b', 'c', 'd', 'e', 'f', 'g'}

	cache := alloc.NewLargeBuffer()
	http := (RequestAuthenticatorFactory{}).Create(&RequestConfig{
		Method:  "GET",
		Uri:     []string{"/"},
		Version: "1.1",
		Header: []*Header{
			{
				Name:  "Content-Length",
				Value: []string{"123"},
			},
		},
	})

	http.Seal(cache).Write(content)

	actualContent := make([]byte, 256)
	reader, err := http.Open(cache)
	assert.Error(err).IsNil()

	n, err := reader.Read(actualContent)
	assert.Error(err).IsNil()
	assert.Bytes(content).Equals(actualContent[:n])
}
