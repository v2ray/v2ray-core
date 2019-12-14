package control

import (
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
)

type FetchCommand struct{}

func (c *FetchCommand) Name() string {
	return "fetch"
}

func (c *FetchCommand) Description() Description {
	return Description{
		Short: "Fetch resources",
		Usage: []string{"v2ctl fetch <url>"},
	}
}

func (c *FetchCommand) Execute(args []string) error {
	if len(args) < 1 {
		return newError("empty url")
	}
	content, err := FetchHTTPContent(args[0])
	if err != nil {
		return newError("failed to read HTTP response").Base(err)
	}

	os.Stdout.Write(content)
	return nil
}

func FetchHTTPContent(target string) ([]byte, error) {

	parsedTarget, err := url.Parse(target)
	if err != nil {
		return nil, newError("invalid URL: ", target).Base(err)
	}

	if s := strings.ToLower(parsedTarget.Scheme); s != "http" && s != "https" {
		return nil, newError("invalid scheme: ", parsedTarget.Scheme)
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Do(&http.Request{
		Method: "GET",
		URL:    parsedTarget,
		Close:  true,
	})
	if err != nil {
		return nil, newError("failed to dial to ", target).Base(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, newError("unexpected HTTP status code: ", resp.StatusCode)
	}

	content, err := buf.ReadAllToBytes(resp.Body)
	if err != nil {
		return nil, newError("failed to read HTTP response").Base(err)
	}

	return content, nil
}

func init() {
	common.Must(RegisterCommand(&FetchCommand{}))
}
