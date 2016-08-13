package ws

import (
	"encoding/json"
)

func (this *Config) UnmarshalJSON(data []byte) error {
	type JsonConfig struct {
		ConnectionReuse bool   `json:"connectionReuse"`
		Path            string `json:"Path"`
		Pto             string `json:"Pto"`
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
	return nil
}
