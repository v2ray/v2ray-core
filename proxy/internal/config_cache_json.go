// +build json

package internal

import (
	"github.com/v2ray/v2ray-core/common/loader"
)

func init() {
	inboundConfigCache = loader.NewJSONConfigLoader("protocol", "settings")
	outboundConfigCache = loader.NewJSONConfigLoader("protocol", "settings")
}
