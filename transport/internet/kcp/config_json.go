// +build json

package kcp

import (
	"encoding/json"

	"v2ray.com/core/common"
	"v2ray.com/core/common/log"
	"v2ray.com/core/transport/internet"
)

func (this *Config) UnmarshalJSON(data []byte) error {
	type JSONConfig struct {
		Mtu             *uint32         `json:"mtu"`
		Tti             *uint32         `json:"tti"`
		UpCap           *uint32         `json:"uplinkCapacity"`
		DownCap         *uint32         `json:"downlinkCapacity"`
		Congestion      *bool           `json:"congestion"`
		ReadBufferSize  *uint32         `json:"readBufferSize"`
		WriteBufferSize *uint32         `json:"writeBufferSize"`
		HeaderConfig    json.RawMessage `json:"header"`
	}
	jsonConfig := new(JSONConfig)
	if err := json.Unmarshal(data, &jsonConfig); err != nil {
		return err
	}
	if jsonConfig.Mtu != nil {
		mtu := *jsonConfig.Mtu
		if mtu < 576 || mtu > 1460 {
			log.Error("KCP|Config: Invalid MTU size: ", mtu)
			return common.ErrBadConfiguration
		}
		this.Mtu = mtu
	}
	if jsonConfig.Tti != nil {
		tti := *jsonConfig.Tti
		if tti < 10 || tti > 100 {
			log.Error("KCP|Config: Invalid TTI: ", tti)
			return common.ErrBadConfiguration
		}
		this.Tti = tti
	}
	if jsonConfig.UpCap != nil {
		this.UplinkCapacity = *jsonConfig.UpCap
	}
	if jsonConfig.DownCap != nil {
		this.DownlinkCapacity = *jsonConfig.DownCap
	}
	if jsonConfig.Congestion != nil {
		this.Congestion = *jsonConfig.Congestion
	}
	if jsonConfig.ReadBufferSize != nil {
		size := *jsonConfig.ReadBufferSize
		if size > 0 {
			this.ReadBuffer = size * 1024 * 1024
		} else {
			this.ReadBuffer = 512 * 1024
		}
	}
	if jsonConfig.WriteBufferSize != nil {
		size := *jsonConfig.WriteBufferSize
		if size > 0 {
			this.WriteBuffer = size * 1024 * 1024
		} else {
			this.WriteBuffer = 512 * 1024
		}
	}
	if len(jsonConfig.HeaderConfig) > 0 {
		name, config, err := internet.CreateAuthenticatorConfig(jsonConfig.HeaderConfig)
		if err != nil {
			log.Error("KCP|Config: Failed to parse header config: ", err)
			return err
		}
		authConfig, err := internet.NewAuthenticatorConfig(name, config)
		if err != nil {
			log.Error("KCP:Config: Failed to create header config: ", err)
			return err
		}
		this.HeaderConfig = authConfig
	}

	return nil
}
