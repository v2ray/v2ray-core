package json

import (
	"encoding/json"

	"github.com/v2ray/v2ray-core/common/log"
)

func JsonConfigLoader(newConfig func() interface{}) func(data []byte) (interface{}, error) {
	return func(data []byte) (interface{}, error) {
		obj := newConfig()
		log.Debug("Unmarshalling JSON: %s", string(data))
		err := json.Unmarshal(data, obj)
		return obj, err
	}
}
