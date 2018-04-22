package jsonem

import (
	"v2ray.com/core"
	"v2ray.com/core/common"
	"v2ray.com/ext/tools/conf/serial"
)

func init() {
	common.Must(core.RegisterConfigLoader(&core.ConfigFormat{
		Name:      "JSON",
		Extension: []string{"json"},
		Loader:    serial.LoadJSONConfig,
	}))
}
