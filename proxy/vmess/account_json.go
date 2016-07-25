// +build json

package vmess

import (
	"encoding/json"

	"github.com/v2ray/v2ray-core/common/log"
	"github.com/v2ray/v2ray-core/common/protocol"
	"github.com/v2ray/v2ray-core/common/uuid"
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
	id, err := uuid.ParseString(rawConfig.ID)
	if err != nil {
		log.Error("VMess: Failed to parse ID: ", err)
		return err
	}
	u.ID = protocol.NewID(id)
	u.AlterIDs = protocol.NewAlterIDs(u.ID, rawConfig.AlterIds)

	return nil
}
