// +build json

package kcp

import (
	"encoding/json"
	"errors"

	"github.com/v2ray/v2ray-core/common/log"
)

func (this *Config) UnmarshalJSON(data []byte) error {
	type JSONConfig struct {
		Mtu *int `json:"mtu"`
	}
	jsonConfig := new(JSONConfig)
	if err := json.Unmarshal(data, &jsonConfig); err != nil {
		return err
	}
	if jsonConfig.Mtu != nil {
		mtu := *jsonConfig.Mtu
		if mtu < 576 || mtu > 1460 {
			log.Error("KCP|Config: Invalid MTU size: ", mtu)
			return errors.New("Invalid configuration")
		}
		this.Mtu = *jsonConfig.Mtu
	}

	return nil
}
