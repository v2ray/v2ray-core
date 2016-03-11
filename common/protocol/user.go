package protocol

import (
	"github.com/v2ray/v2ray-core/common/dice"
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
	Email    string
}

func NewUser(id *ID, level UserLevel, alterIdCount uint16, email string) *User {
	u := &User{
		ID:    id,
		Level: level,
		Email: email,
	}
	if alterIdCount > 0 {
		u.AlterIDs = make([]*ID, alterIdCount)
		prevId := id.UUID()
		for idx := range u.AlterIDs {
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
	return this.AlterIDs[dice.Roll(len(this.AlterIDs))]
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

type Account interface {
	CryptionKey() []byte
}
