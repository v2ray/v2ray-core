package control

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/golang/protobuf/proto"
	"v2ray.com/core/common"
	"v2ray.com/core/infra/conf"
	"v2ray.com/core/infra/conf/serial"
)

type MconfigCommand struct{}

func (c *MconfigCommand) Name() string {
	return "mconfig"
}

func (c *MconfigCommand) Description() Description {
	return Description{
		Short: "merge multiple json config",
		Usage: []string{"v2ctl mconfig 1.json 2.json <url>.json"},
	}
}

func (c *MconfigCommand) Execute(args []string) error {
	if len(args) < 1 {
		return newError("empty config list")
	}

	conf := &conf.Config{}
	for _, arg := range args {
		r, err := c.LoadArg(arg)
		common.Must(err)
		c, err := serial.DecodeJSONConfig(r)
		common.Must(err)
		conf.Override(c, arg)
	}

	pbConfig, err := conf.Build()
	if err != nil {
		return err
	}

	bytesConfig, err := proto.Marshal(pbConfig)
	if err != nil {
		return newError("failed to marshal proto config").Base(err)
	}

	if _, err := os.Stdout.Write(bytesConfig); err != nil {
		return newError("failed to write proto config").Base(err)
	}

	return nil
}

func (c *MconfigCommand) LoadArg(arg string) (out io.Reader, err error) {

	var data []byte
	if strings.HasPrefix(arg, "http://") || strings.HasPrefix(arg, "https://") {
		data, err = FetchHTTPContent(arg)
	} else {
		data, err = ioutil.ReadFile(arg)
	}

	if err != nil {
		return
	}
	out = bytes.NewBuffer(data)
	return
}

func init() {
	common.Must(RegisterCommand(&MconfigCommand{}))
}
