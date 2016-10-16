// +build json

package internet

import (
	"encoding/json"
	"strings"

	"v2ray.com/core/common/loader"
	v2net "v2ray.com/core/common/net"
	v2tls "v2ray.com/core/transport/internet/tls"
)

func (this *StreamConfig) UnmarshalJSON(data []byte) error {
	type JSONConfig struct {
		Network     *v2net.Network `json:"network"`
		Security    string         `json:"security"`
		TLSSettings *v2tls.Config  `json:"tlsSettings"`
	}
	this.Network = v2net.Network_RawTCP
	jsonConfig := new(JSONConfig)
	if err := json.Unmarshal(data, jsonConfig); err != nil {
		return err
	}
	if jsonConfig.Network != nil {
		this.Network = *jsonConfig.Network
	}
	if strings.ToLower(jsonConfig.Security) == "tls" {
		tlsSettings := jsonConfig.TLSSettings
		if tlsSettings == nil {
			tlsSettings = &v2tls.Config{}
		}
		this.SecuritySettings = append(this.SecuritySettings, loader.NewTypedSettings(tlsSettings))
	}
	return nil
}
