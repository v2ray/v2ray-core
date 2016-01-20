package vmess

import (
	"math/rand"
)

type UserLevel byte

const (
	UserLevelAdmin     = UserLevel(255)
	UserLevelUntrusted = UserLevel(0)
)

type User struct {
	ID       *ID
	AlterIDs []*ID
	Level    UserLevel
}

func NewUser(id *ID, level UserLevel, alterIdCount uint16) *User {
	u := &User{
		ID:    id,
		Level: level,
	}
	if alterIdCount > 0 {
		u.AlterIDs = make([]*ID, alterIdCount)
		prevId := id.UUID()
		for idx, _ := range u.AlterIDs {
			newid := prevId.Next()
			// TODO: check duplicate
			u.AlterIDs[idx] = NewID(newid)
			prevId = newid
		}
	}
	return u
}

func (this *User) AnyValidID() *ID {
	if len(this.AlterIDs) == 0 {
		return this.ID
	}
	if len(this.AlterIDs) == 1 {
		return this.AlterIDs[0]
	}
	idx := rand.Intn(len(this.AlterIDs))
	return this.AlterIDs[idx]
}

type UserSettings struct {
	PayloadReadTimeout int
}

func GetUserSettings(level UserLevel) UserSettings {
	settings := UserSettings{
		PayloadReadTimeout: 120,
	}
	if level > 0 {
		settings.PayloadReadTimeout = 0
	}
	return settings
}
