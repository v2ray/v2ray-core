package json

//go:generate errorgen

import (
	"io"
	"os"

	"v2ray.com/core"
	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/cmdarg"
	"v2ray.com/core/common/platform/ctlcmd"
	"v2ray.com/core/infra/conf/serial"
)

func init() {
	common.Must(core.RegisterConfigLoader(&core.ConfigFormat{
		Name:      "JSON",
		Extension: []string{"json"},
		Loader: func(input interface{}) (*core.Config, error) {
			switch v := input.(type) {
			case cmdarg.Arg:
				jsonContent, err := ctlcmd.Run(append([]string{"config"}, v...), os.Stdin)
				if err != nil {
					return nil, newError("failed to execute v2ctl to convert config file.").Base(err).AtWarning()
				}
				return core.LoadConfig("protobuf", "", &buf.MultiBufferContainer{
					MultiBuffer: jsonContent,
				})
			case io.Reader:
				return serial.LoadJSONConfig(v)
			default:
				return nil, newError("unknow type")
			}
		},
	}))
}
