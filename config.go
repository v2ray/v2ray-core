package core

import (
	"encoding/json"
)

// User is the user account that is used for connection to a Point
type User struct {
	Id ID `json:"id"` // The ID of this User.
}

type ConnectionConfig struct {
	Protocol string `json:"protocol"`
	File     string `json:"file"`
}

// Config is the config for Point server.
type Config struct {
	Port           uint16            `json:"port"` // Port of this Point server.
	InboundConfig  ConnectionConfig `json:"inbound"`
	OutboundConfig ConnectionConfig `json:"outbound"`
}

func LoadConfig(rawConfig []byte) (Config, error) {
	config := Config{}
	err := json.Unmarshal(rawConfig, &config)
	return config, err
}
