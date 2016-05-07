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

func NewUser(primary *ID, secondary []*ID, level UserLevel, email string) *User {
	return &User{
		ID:       primary,
		AlterIDs: secondary,
		Level:    level,
		Email:    email,
	}
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
