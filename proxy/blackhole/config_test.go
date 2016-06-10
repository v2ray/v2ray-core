package blackhole_test

import (
	"bufio"
	"net/http"
	"testing"

	"github.com/v2ray/v2ray-core/common/alloc"
	v2io "github.com/v2ray/v2ray-core/common/io"
	. "github.com/v2ray/v2ray-core/proxy/blackhole"
	"github.com/v2ray/v2ray-core/testing/assert"
)

func TestHTTPResponse(t *testing.T) {
	assert := assert.On(t)

	buffer := alloc.NewBuffer().Clear()

	httpResponse := new(HTTPResponse)
	httpResponse.WriteTo(v2io.NewAdaptiveWriter(buffer))

	reader := bufio.NewReader(buffer)
	response, err := http.ReadResponse(reader, nil)
	assert.Error(err).IsNil()
	assert.Int(response.StatusCode).Equals(403)
}
