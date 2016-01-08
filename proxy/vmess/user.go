package vmess

type UserLevel int

const (
	UserLevelAdmin     = UserLevel(999)
	UserLevelUntrusted = UserLevel(0)
)

type User interface {
	ID() *ID
	AlterIDs() []*ID
	Level() UserLevel
	AnyValidID() *ID
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
