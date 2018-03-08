package blackhole_test

import (
	"bufio"
	"net/http"
	"testing"

	"v2ray.com/core/common/buf"
	. "v2ray.com/core/proxy/blackhole"
	. "v2ray.com/ext/assert"
)

func TestHTTPResponse(t *testing.T) {
	assert := With(t)

	buffer := buf.New()

	httpResponse := new(HTTPResponse)
	httpResponse.WriteTo(buf.NewWriter(buffer))

	reader := bufio.NewReader(buffer)
	response, err := http.ReadResponse(reader, nil)
	assert(err, IsNil)
	assert(response.StatusCode, Equals, 403)
}
