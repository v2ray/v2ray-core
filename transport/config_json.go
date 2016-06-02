// +build json

package transport

import (
	"encoding/json"
	"strings"
)

func (this *Config) UnmarshalJSON(data []byte) error {
	type TypeConfig struct {
		StreamType string          `json:"streamType"`
		Settings   json.RawMessage `json:"settings"`
	}
	type JsonTCPConfig struct {
		ConnectionReuse bool `json:"connectionReuse"`
	}

	typeConfig := new(TypeConfig)
	if err := json.Unmarshal(data, typeConfig); err != nil {
		return err
	}

	this.StreamType = StreamTypeTCP

	streamType := strings.ToLower(typeConfig.StreamType)
	if streamType == "tcp" {
		jsonTCPConfig := new(JsonTCPConfig)
		if err := json.Unmarshal(typeConfig.Settings, jsonTCPConfig); err != nil {
			return err
		}
		this.TCPConfig = &TCPConfig{
			ConnectionReuse: jsonTCPConfig.ConnectionReuse,
		}
	}

	return nil
}
