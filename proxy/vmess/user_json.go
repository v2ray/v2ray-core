// +build json

package vmess

import (
	"encoding/json"

	"github.com/v2ray/v2ray-core/common/uuid"
)

func (u *User) UnmarshalJSON(data []byte) error {
	type rawUser struct {
		IdString     string `json:"id"`
		EmailString  string `json:"email"`
		LevelInt     int    `json:"level"`
		AlterIdCount int    `json:"alterId"`
	}
	var rawUserValue rawUser
	if err := json.Unmarshal(data, &rawUserValue); err != nil {
		return err
	}
	id, err := uuid.ParseString(rawUserValue.IdString)
	if err != nil {
		return err
	}
	u.ID = NewID(id)
	//u.Email = rawUserValue.EmailString
	u.Level = UserLevel(rawUserValue.LevelInt)

	if rawUserValue.AlterIdCount > 0 {
		prevId := u.ID.UUID()
		// TODO: check duplicate
		u.AlterIDs = make([]*ID, rawUserValue.AlterIdCount)
		for idx, _ := range u.AlterIDs {
			newid := prevId.Next()
			u.AlterIDs[idx] = NewID(newid)
			prevId = newid
		}
	}

	return nil
}
