package scenarios

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"

	"v2ray.com/core/common/net"
	"v2ray.com/core/testing/assert"
	v2http "v2ray.com/core/testing/servers/http"
)

func TestHttpProxy(t *testing.T) {
	assert := assert.On(t)

	httpServer := &v2http.Server{
		Port:        net.Port(50042),
		PathHandler: make(map[string]http.HandlerFunc),
	}
	_, err := httpServer.Start()
	assert.Error(err).IsNil()
	defer httpServer.Close()

	assert.Error(InitializeServerSetOnce("test_5")).IsNil()

	transport := &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return url.Parse("http://127.0.0.1:50040/")
		},
	}

	client := &http.Client{
		Transport: transport,
	}

	resp, err := client.Get("http://127.0.0.1:50042/")
	assert.Error(err).IsNil()
	assert.Int(resp.StatusCode).Equals(200)

	content, err := ioutil.ReadAll(resp.Body)
	assert.Error(err).IsNil()
	assert.String(string(content)).Equals("Home")

	CloseAllServers()
}

func TestBlockHTTP(t *testing.T) {
	assert := assert.On(t)

	httpServer := &v2http.Server{
		Port:        net.Port(50042),
		PathHandler: make(map[string]http.HandlerFunc),
	}
	_, err := httpServer.Start()
	assert.Error(err).IsNil()
	defer httpServer.Close()

	assert.Error(InitializeServerSetOnce("test_5")).IsNil()

	transport := &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return url.Parse("http://127.0.0.1:50040/")
		},
	}

	client := &http.Client{
		Transport: transport,
	}

	resp, err := client.Get("http://127.0.0.1:50049/")
	assert.Error(err).IsNil()
	assert.Int(resp.StatusCode).Equals(403)

	CloseAllServers()
}
