package conf

import (
	"io"

	"v2ray.com/core"
	jsonconf "v2ray.com/ext/tools/conf/serial"
)

//go:generate go run $GOPATH/src/v2ray.com/core/tools/generrorgen/main.go -pkg conf -path Tools,Conf

func init() {
	core.RegisterConfigLoader(core.ConfigFormat_JSON, func(input io.Reader) (*core.Config, error) {
		return jsonconf.LoadJSONConfig(input)
	})
}
