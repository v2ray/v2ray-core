// +build json

package http

import (
	"crypto/tls"
	"encoding/json"

	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/proxy/internal/config"
)

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

func (this *TlsConfig) UnmarshalJSON(data []byte) error {
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

func (this *Config) UnmarshalJSON(data []byte) error {
	type JsonConfig struct {
		Hosts []v2net.AddressJson `json:"ownHosts"`
		Tls   *TlsConfig          `json:"tls"`
	}
	jsonConfig := new(JsonConfig)
	if err := json.Unmarshal(data, jsonConfig); err != nil {
		return err
	}
	this.OwnHosts = make([]v2net.Address, len(jsonConfig.Hosts), len(jsonConfig.Hosts)+1)
	for idx, host := range jsonConfig.Hosts {
		this.OwnHosts[idx] = host.Address
	}

	v2rayHost := v2net.DomainAddress("local.v2ray.com")
	if !this.IsOwnHost(v2rayHost) {
		this.OwnHosts = append(this.OwnHosts, v2rayHost)
	}

	this.TlsConfig = jsonConfig.Tls

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
