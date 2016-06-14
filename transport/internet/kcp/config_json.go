// +build json

package kcp

import (
	"encoding/json"
)

func (this *Config) UnmarshalJSON(data []byte) error {
	type JSONConfig struct {
		Mode         string `json:"Mode"`
		Mtu          int    `json:"MaximumTransmissionUnit"`
		Sndwnd       int    `json:"SendingWindowSize"`
		Rcvwnd       int    `json:"ReceivingWindowSize"`
		Fec          int    `json:"ForwardErrorCorrectionGroupSize"`
		Acknodelay   bool   `json:"AcknowledgeNoDelay"`
		Dscp         int    `json:"Dscp"`
		ReadTimeout  int    `json:"ReadTimeout"`
		WriteTimeout int    `json:"WriteTimeout"`
	}
	jsonConfig := effectiveConfig
	if err := json.Unmarshal(data, &jsonConfig); err != nil {
		return err
	}
	*this = jsonConfig
	return nil
}
