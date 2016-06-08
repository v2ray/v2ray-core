// +build json

package http

import (
	"crypto/tls"
	"encoding/json"

	"github.com/v2ray/v2ray-core/proxy/internal/config"
)

// UnmarshalJSON implements json.Unmarshaler
func (this *CertificateConfig) UnmarshalJSON(data []byte) error {
	type JsonConfig struct {
		Domain   string `json:"domain"`
		CertFile string `json:"cert"`
		KeyFile  string `json:"key"`
	}
	jsonConfig := new(JsonConfig)
	if err := json.Unmarshal(data, jsonConfig); err != nil {
		return err
	}

	cert, err := tls.LoadX509KeyPair(jsonConfig.CertFile, jsonConfig.KeyFile)
	if err != nil {
		return err
	}
	this.Domain = jsonConfig.Domain
	this.Certificate = cert
	return nil
}

// UnmarshalJSON implements json.Unmarshaler
func (this *TLSConfig) UnmarshalJSON(data []byte) error {
	type JsonConfig struct {
		Enabled bool                 `json:"enable"`
		Certs   []*CertificateConfig `json:"certs"`
	}
	jsonConfig := new(JsonConfig)
	if err := json.Unmarshal(data, jsonConfig); err != nil {
		return err
	}

	this.Enabled = jsonConfig.Enabled
	this.Certs = jsonConfig.Certs
	return nil
}

// UnmarshalJSON implements json.Unmarshaler
func (this *Config) UnmarshalJSON(data []byte) error {
	type JsonConfig struct {
		Tls *TLSConfig `json:"tls"`
	}
	jsonConfig := new(JsonConfig)
	if err := json.Unmarshal(data, jsonConfig); err != nil {
		return err
	}

	this.TLSConfig = jsonConfig.Tls

	return nil
}

func init() {
	config.RegisterInboundConfig("http",
		func(data []byte) (interface{}, error) {
			rawConfig := new(Config)
			err := json.Unmarshal(data, rawConfig)
			return rawConfig, err
		})
}
