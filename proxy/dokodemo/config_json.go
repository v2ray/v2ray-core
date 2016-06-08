// +build json

package dokodemo

import (
	"encoding/json"

	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/proxy/internal/config"
)

func init() {
	config.RegisterInboundConfig("dokodemo-door",
		func(data []byte) (interface{}, error) {
			type DokodemoConfig struct {
				Host         *v2net.AddressJson `json:"address"`
				PortValue    v2net.Port         `json:"port"`
				NetworkList  *v2net.NetworkList `json:"network"`
				TimeoutValue int                `json:"timeout"`
			}
			rawConfig := new(DokodemoConfig)
			if err := json.Unmarshal(data, rawConfig); err != nil {
				return nil, err
			}
			return &Config{
				Address: rawConfig.Host.Address,
				Port:    rawConfig.PortValue,
				Network: rawConfig.NetworkList,
				Timeout: rawConfig.TimeoutValue,
			}, nil
		})
}
