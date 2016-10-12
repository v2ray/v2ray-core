// +build json

package vmess

import (
	"encoding/json"
)

func (u *Account) UnmarshalJSON(data []byte) error {
	type JsonConfig struct {
		ID       string `json:"id"`
		AlterIds uint16 `json:"alterId"`
	}
	var rawConfig JsonConfig
	if err := json.Unmarshal(data, &rawConfig); err != nil {
		return err
	}
	u.Id = rawConfig.ID
	u.AlterId = uint32(rawConfig.AlterIds)

	return nil
}
