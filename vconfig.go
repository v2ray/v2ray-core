package core

import (
	"encoding/json"
)

// VUser is the user account that is used for connection to a VPoint
type VUser struct {
	Id VID `json:"id"` // The ID of this VUser.
}

type VConnectionConfig struct {
	Protocol string `json:"protocol"`
	File     string `json:"file"`
}

// VConfig is the config for VPoint server.
type VConfig struct {
	Port           uint16            `json:"port"` // Port of this VPoint server.
	InboundConfig  VConnectionConfig `json:"inbound"`
	OutboundConfig VConnectionConfig `json:"outbound"`
}

func LoadVConfig(rawConfig []byte) (VConfig, error) {
	config := VConfig{}
	err := json.Unmarshal(rawConfig, &config)
	return config, err
}
