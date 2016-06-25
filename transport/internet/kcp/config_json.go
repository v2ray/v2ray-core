// +build json

package kcp

import (
	"encoding/json"

	"github.com/v2ray/v2ray-core/common"
	"github.com/v2ray/v2ray-core/common/log"
)

func (this *Config) UnmarshalJSON(data []byte) error {
	type JSONConfig struct {
		Mtu        *uint32 `json:"mtu"`
		Tti        *uint32 `json:"tti"`
		UpCap      *uint32 `json:"uplinkCapacity"`
		DownCap    *uint32 `json:"downlinkCapacity"`
		Congestion *bool   `json:"congestion"`
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
		upCap := *jsonConfig.UpCap
		if upCap == 0 {
			log.Error("KCP|Config: Invalid uplink capacity: ", upCap)
			return common.ErrBadConfiguration
		}
		this.UplinkCapacity = upCap
	}
	if jsonConfig.DownCap != nil {
		downCap := *jsonConfig.DownCap
		if downCap == 0 {
			log.Error("KCP|Config: Invalid downlink capacity: ", downCap)
			return common.ErrBadConfiguration
		}
		this.DownlinkCapacity = downCap
	}
	if jsonConfig.Congestion != nil {
		this.Congestion = *jsonConfig.Congestion
	}

	return nil
}
