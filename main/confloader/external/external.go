package external

//go:generate go run v2ray.com/core/common/errors/errorgen

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/platform/ctlcmd"
	"v2ray.com/core/main/confloader"
)

func ConfigLoader(arg string) (out io.Reader, err error) {

	var data []byte
	if strings.HasPrefix(arg, "http://") || strings.HasPrefix(arg, "https://") {
		data, err = FetchHTTPContent(arg)
	} else if arg == "stdin:" {
		data, err = ioutil.ReadAll(os.Stdin)
	} else {
		data, err = ioutil.ReadFile(arg)
	}

	if err != nil {
		return
	}
	out = bytes.NewBuffer(data)
	return
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

func ExtConfigLoader(files []string) (io.Reader, error) {
	buf, err := ctlcmd.Run(append([]string{"config"}, files...), os.Stdin)
	if err != nil {
		return nil, err
	}

	return strings.NewReader(buf.String()), nil
}

func init() {
	confloader.EffectiveConfigFileLoader = ConfigLoader
	confloader.EffectiveExtConfigLoader = ExtConfigLoader
}
