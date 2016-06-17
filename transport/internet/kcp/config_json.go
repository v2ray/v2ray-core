// +build json

package kcp

import (
	"encoding/json"
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
		this.Mtu = *jsonConfig.Mtu
	}

	return nil
}
