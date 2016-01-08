package json

import (
	"encoding/json"
	"math/rand"

	"github.com/v2ray/v2ray-core/common/uuid"
	"github.com/v2ray/v2ray-core/proxy/vmess"
)

// ConfigUser is an user account in VMess configuration.
type ConfigUser struct {
	Id         *vmess.ID
	Email      string
	LevelValue vmess.UserLevel
	AlterIds   []*vmess.ID
}

func (u *ConfigUser) UnmarshalJSON(data []byte) error {
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
	u.Id = vmess.NewID(id)
	u.Email = rawUserValue.EmailString
	u.LevelValue = vmess.UserLevel(rawUserValue.LevelInt)

	if rawUserValue.AlterIdCount > 0 {
		prevId := u.Id.UUID()
		// TODO: check duplicate
		u.AlterIds = make([]*vmess.ID, rawUserValue.AlterIdCount)
		for idx, _ := range u.AlterIds {
			newid := prevId.Next()
			u.AlterIds[idx] = vmess.NewID(newid)
			prevId = newid
		}
	}

	return nil
}

func (u *ConfigUser) ID() *vmess.ID {
	return u.Id
}

func (this *ConfigUser) Level() vmess.UserLevel {
	return this.LevelValue
}

func (this *ConfigUser) AlterIDs() []*vmess.ID {
	return this.AlterIds
}

func (this *ConfigUser) AnyValidID() *vmess.ID {
	if len(this.AlterIds) == 0 {
		return this.ID()
	}
	if len(this.AlterIds) == 1 {
		return this.AlterIds[0]
	}
	idIdx := rand.Intn(len(this.AlterIds))
	return this.AlterIds[idIdx]
}
