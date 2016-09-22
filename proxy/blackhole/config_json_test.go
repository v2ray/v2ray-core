// +build json

package blackhole_test

import (
	"encoding/json"
	"testing"

	. "v2ray.com/core/proxy/blackhole"
	"v2ray.com/core/testing/assert"
)

func TestHTTPResponseJSON(t *testing.T) {
	assert := assert.On(t)

	rawJson := `{
    "response": {
      "type": "http"
    }
  }`
	config := new(Config)
	err := json.Unmarshal([]byte(rawJson), config)
	assert.Error(err).IsNil()

	assert.Int(int(config.Response.Type)).Equals(int(Response_HTTP))
	response, err := config.Response.GetInternalResponse()
	assert.Error(err).IsNil()

	_, ok := response.(*HTTPResponse)
	assert.Bool(ok).IsTrue()
}
