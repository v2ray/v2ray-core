// +build json

package internet

import (
	"encoding/json"
	"strings"

	"errors"
	"github.com/golang/protobuf/ptypes"
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
	this.SecurityType = SecurityType_None
	if strings.ToLower(jsonConfig.Security) == "tls" {
		this.SecurityType = SecurityType_TLS
	}
	if jsonConfig.TLSSettings != nil {
		anyTLSSettings, err := ptypes.MarshalAny(jsonConfig.TLSSettings)
		if err != nil {
			return errors.New("Internet: Failed to parse TLS settings: " + err.Error())
		}
		this.SecuritySettings = append(this.SecuritySettings, &SecuritySettings{
			Type:     SecurityType_TLS,
			Settings: anyTLSSettings,
		})
	}
	return nil
}
