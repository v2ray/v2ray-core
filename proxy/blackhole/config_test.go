package blackhole_test

import (
	"bufio"
	"net/http"
	"testing"

	"v2ray.com/core/common/buf"
	v2io "v2ray.com/core/common/io"
	. "v2ray.com/core/proxy/blackhole"
	"v2ray.com/core/testing/assert"
)

func TestHTTPResponse(t *testing.T) {
	assert := assert.On(t)

	buffer := buf.NewBuffer()

	httpResponse := new(HTTPResponse)
	httpResponse.WriteTo(v2io.NewAdaptiveWriter(buffer))

	reader := bufio.NewReader(buffer)
	response, err := http.ReadResponse(reader, nil)
	assert.Error(err).IsNil()
	assert.Int(response.StatusCode).Equals(403)
}
