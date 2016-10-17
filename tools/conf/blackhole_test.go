package conf_test

import (
	"encoding/json"
	"testing"

	"v2ray.com/core/proxy/blackhole"
	"v2ray.com/core/testing/assert"
	. "v2ray.com/core/tools/conf"
)

func TestHTTPResponseJSON(t *testing.T) {
	assert := assert.On(t)

	rawJson := `{
    "response": {
      "type": "http"
    }
  }`
	rawConfig := new(BlackholeConfig)
	err := json.Unmarshal([]byte(rawJson), rawConfig)
	assert.Error(err).IsNil()

	ts, err := rawConfig.Build()
	assert.Error(err).IsNil()
	iConfig, err := ts.GetInstance()
	assert.Error(err).IsNil()
	config := iConfig.(*blackhole.Config)
	response, err := config.GetInternalResponse()
	assert.Error(err).IsNil()

	_, ok := response.(*blackhole.HTTPResponse)
	assert.Bool(ok).IsTrue()
}
