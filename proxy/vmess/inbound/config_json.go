// +build json

package inbound

import (
	"encoding/json"

	"github.com/v2ray/v2ray-core/proxy/internal/config"
	"github.com/v2ray/v2ray-core/proxy/vmess"
)

func init() {
	config.RegisterInboundConnectionConfig("vmess",
		func(data []byte) (interface{}, error) {
			type JsonConfig struct {
				Users []*vmess.User `json:"clients"`
			}
			jsonConfig := new(JsonConfig)
			if err := json.Unmarshal(data, jsonConfig); err != nil {
				return nil, err
			}
			return &Config{
				AllowedUsers: jsonConfig.Users,
			}, nil
		})
}
