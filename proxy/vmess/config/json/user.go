package json

import (
	"encoding/json"

	"github.com/v2ray/v2ray-core/proxy/vmess/config"
)

// ConfigUser is an user account in VMess configuration.
type ConfigUser struct {
	Id    *config.ID
	Email string
  LevelValue config.UserLevel
}

func (u *ConfigUser) UnmarshalJSON(data []byte) error {
	type rawUser struct {
		IdString    string `json:"id"`
		EmailString string `json:"email"`
    LevelInt int `json:"level"`
	}
	var rawUserValue rawUser
	if err := json.Unmarshal(data, &rawUserValue); err != nil {
		return err
	}
	id, err := config.NewID(rawUserValue.IdString)
	if err != nil {
		return err
	}
	u.Id = id
	u.Email = rawUserValue.EmailString
  u.LevelValue = config.UserLevel(rawUserValue.LevelInt)
	return nil
}

func (u *ConfigUser) ID() *config.ID {
	return u.Id
}

func (this *ConfigUser) Level() config.UserLevel {
  return this.LevelValue
}