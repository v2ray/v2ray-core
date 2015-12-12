package json

import (
	"encoding/json"

	"github.com/v2ray/v2ray-core/common/uuid"
	"github.com/v2ray/v2ray-core/proxy/vmess"
)

// ConfigUser is an user account in VMess configuration.
type ConfigUser struct {
	Id         *vmess.ID
	Email      string
	LevelValue vmess.UserLevel
}

func (u *ConfigUser) UnmarshalJSON(data []byte) error {
	type rawUser struct {
		IdString    string `json:"id"`
		EmailString string `json:"email"`
		LevelInt    int    `json:"level"`
	}
	var rawUserValue rawUser
	if err := json.Unmarshal(data, &rawUserValue); err != nil {
		return err
	}
	id, err := uuid.ParseString(rawUserValue.IdString)
	if err != nil {
		return err
	}
	u.Id = vmess.NewID(id)
	u.Email = rawUserValue.EmailString
	u.LevelValue = vmess.UserLevel(rawUserValue.LevelInt)
	return nil
}

func (u *ConfigUser) ID() *vmess.ID {
	return u.Id
}

func (this *ConfigUser) Level() vmess.UserLevel {
	return this.LevelValue
}
