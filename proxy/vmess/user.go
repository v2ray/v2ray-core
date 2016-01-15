package vmess

import (
	"math/rand"
)

type UserLevel int

const (
	UserLevelAdmin     = UserLevel(999)
	UserLevelUntrusted = UserLevel(0)
)

type User struct {
	ID       *ID
	AlterIDs []*ID
	Level    UserLevel
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
