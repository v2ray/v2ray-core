// +build json

package internet

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"strings"

	v2net "github.com/v2ray/v2ray-core/common/net"
)

func (this *TLSSettings) UnmarshalJSON(data []byte) error {
	type JSONCertConfig struct {
		CertFile string `json:"certificateFile"`
		KeyFile  string `json:"keyFile"`
	}
	type JSONConfig struct {
		Certs []*JSONCertConfig `json:"certificates"`
	}
	jsonConfig := new(JSONConfig)
	if err := json.Unmarshal(data, jsonConfig); err != nil {
		return err
	}
	this.Certs = make([]tls.Certificate, len(jsonConfig.Certs))
	for idx, certConf := range jsonConfig.Certs {
		cert, err := tls.LoadX509KeyPair(certConf.CertFile, certConf.KeyFile)
		if err != nil {
			return errors.New("Internet|TLS: Failed to load certificate file: " + err.Error())
		}
		this.Certs[idx] = cert
	}
	return nil
}

func (this *StreamSettings) UnmarshalJSON(data []byte) error {
	type JSONConfig struct {
		Network     v2net.NetworkList `json:"network"`
		Security    string            `json:"security"`
		TLSSettings *TLSSettings      `json:"tlsSettings"`
	}
	this.Type = StreamConnectionTypeRawTCP
	jsonConfig := new(JSONConfig)
	if err := json.Unmarshal(data, jsonConfig); err != nil {
		return err
	}
	if jsonConfig.Network.HasNetwork(v2net.KCPNetwork) {
		this.Type |= StreamConnectionTypeKCP
	}
	if jsonConfig.Network.HasNetwork(v2net.TCPNetwork) {
		this.Type |= StreamConnectionTypeTCP
	}
	this.Security = StreamSecurityTypeNone
	if strings.ToLower(jsonConfig.Security) == "tls" {
		this.Security = StreamSecurityTypeTLS
	}
	if jsonConfig.TLSSettings != nil {
		this.TLSSettings = jsonConfig.TLSSettings
	}
	return nil
}
