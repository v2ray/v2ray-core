// +build json

package internet

import "github.com/v2ray/v2ray-core/common/loader"

func init() {
	configCache = loader.NewJSONConfigLoader("type", "")
}
