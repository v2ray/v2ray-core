package control

import (
	"net/http"
	"net/url"
	"os"
	"strings"

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

func (c *FetchCommand) isValidScheme(scheme string) bool {
	scheme = strings.ToLower(scheme)
	return scheme == "http" || scheme == "https"
}

func (c *FetchCommand) Execute(args []string) error {
	if len(args) < 1 {
		return newError("empty url")
	}
	target := args[0]
	parsedTarget, err := url.Parse(target)
	if err != nil {
		return newError("invalid URL: ", target).Base(err)
	}
	if !c.isValidScheme(parsedTarget.Scheme) {
		return newError("invalid scheme: ", parsedTarget.Scheme)
	}

	client := &http.Client{}
	resp, err := client.Do(&http.Request{
		Method: "GET",
		URL:    parsedTarget,
		Close:  true,
	})
	if err != nil {
		return newError("failed to dial to ", target).Base(err)
	}

	if resp.StatusCode != 200 {
		return newError("unexpected HTTP status code: ", resp.StatusCode)
	}

	content, err := buf.ReadAllToBytes(resp.Body)
	if err != nil {
		return newError("failed to read HTTP response").Base(err)
	}

	os.Stdout.Write(content)

	return nil
}

func init() {
	common.Must(RegisterCommand(&FetchCommand{}))
}
