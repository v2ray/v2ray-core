package json

//go:generate errorgen

import (
	"encoding/json"
	"io"
	"io/ioutil"

	"v2ray.com/core"
	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/platform/ctlcmd"
)

func init() {
	common.Must(core.RegisterConfigLoader(&core.ConfigFormat{
		Name:      "JSON",
		Extension: []string{"json"},
		Loader: func(input io.Reader) (*core.Config, error) {
			fns := []string{}
			data, _ := ioutil.ReadAll(input)
			json.Unmarshal(data, &fns)
			jsonContent, err := ctlcmd.Run(append([]string{"config"}, fns...), nil)
			if err != nil {
				return nil, newError("failed to execute v2ctl to convert config file.").Base(err).AtWarning()
			}
			return core.LoadConfig("protobuf", "", &buf.MultiBufferContainer{
				MultiBuffer: jsonContent,
			})
		},
	}))
}
