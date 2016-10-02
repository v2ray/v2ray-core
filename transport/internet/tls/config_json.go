// +build json

package tls

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

func (this *Config) UnmarshalJSON(data []byte) error {
	type JSONCertConfig struct {
		CertFile string `json:"certificateFile"`
		KeyFile  string `json:"keyFile"`
	}
	type JSONConfig struct {
		Insecure bool              `json:"allowInsecure"`
		Certs    []*JSONCertConfig `json:"certificates"`
	}
	jsonConfig := new(JSONConfig)
	if err := json.Unmarshal(data, jsonConfig); err != nil {
		return err
	}
	this.Certificate = make([]*Certificate, len(jsonConfig.Certs))
	for idx, certConf := range jsonConfig.Certs {
		cert, err := ioutil.ReadFile(certConf.CertFile)
		if err != nil {
			return errors.New("TLS: Failed to load certificate file: " + err.Error())
		}
		key, err := ioutil.ReadFile(certConf.KeyFile)
		if err != nil {
			return errors.New("TLS: Failed to load key file: " + err.Error())
		}
		this.Certificate[idx] = &Certificate{
			Key:         key,
			Certificate: cert,
		}
	}
	this.AllowInsecure = jsonConfig.Insecure
	return nil
}
