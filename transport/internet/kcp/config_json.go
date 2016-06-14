// +build json

package kcp

import (
	"encoding/json"
)

func (this *Config) UnmarshalJSON(data []byte) error {
	type JSONConfig struct {
		Mode         *string `json:"Mode"`
		Mtu          *int    `json:"MaximumTransmissionUnit"`
		Sndwnd       *int    `json:"SendingWindowSize"`
		Rcvwnd       *int    `json:"ReceivingWindowSize"`
		Acknodelay   *bool   `json:"AcknowledgeNoDelay"`
		Dscp         *int    `json:"Dscp"`
		ReadTimeout  *int    `json:"ReadTimeout"`
		WriteTimeout *int    `json:"WriteTimeout"`
	}
	jsonConfig := new(JSONConfig)
	if err := json.Unmarshal(data, &jsonConfig); err != nil {
		return err
	}
	if jsonConfig.Mode != nil {
		this.Mode = *jsonConfig.Mode
	}

	if jsonConfig.Mtu != nil {
		this.Mtu = *jsonConfig.Mtu
	}

	if jsonConfig.Sndwnd != nil {
		this.Sndwnd = *jsonConfig.Sndwnd
	}
	if jsonConfig.Rcvwnd != nil {
		this.Rcvwnd = *jsonConfig.Rcvwnd
	}
	if jsonConfig.Acknodelay != nil {
		this.Acknodelay = *jsonConfig.Acknodelay
	}
	if jsonConfig.Dscp != nil {
		this.Dscp = *jsonConfig.Dscp
	}
	if jsonConfig.ReadTimeout != nil {
		this.ReadTimeout = *jsonConfig.ReadTimeout
	}
	if jsonConfig.WriteTimeout != nil {
		this.WriteTimeout = *jsonConfig.WriteTimeout
	}

	return nil
}
