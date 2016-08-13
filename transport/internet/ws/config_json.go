package ws

import (
	"encoding/json"
)

func (this *Config) UnmarshalJSON(data []byte) error {
	type JsonConfig struct {
		ConnectionReuse bool   `json:"connectionReuse"`
		Path            string `json:"Path"`
		Pto             string `json:"Pto"`
		Cert            string `json:"Cert"`
		PrivKey         string `json:"PrivKet"`
	}
	jsonConfig := &JsonConfig{
		ConnectionReuse: true,
		Path:            "",
		Pto:             "",
	}
	if err := json.Unmarshal(data, jsonConfig); err != nil {
		return err
	}
	this.ConnectionReuse = jsonConfig.ConnectionReuse
	this.Path = jsonConfig.Path
	this.Pto = jsonConfig.Pto
	this.PrivKey = jsonConfig.PrivKey
	this.Cert = jsonConfig.Cert
	return nil
}
