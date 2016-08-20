// +build json

package internet

import "v2ray.com/core/common/loader"

func init() {
	configCache = loader.NewJSONConfigLoader("type", "")
}
