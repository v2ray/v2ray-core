package protocol

type UserLevel byte

const (
	UserLevelAdmin     = UserLevel(255)
	UserLevelUntrusted = UserLevel(0)
)

type User struct {
	Account Account
	Level   UserLevel
	Email   string
}

func NewUser(level UserLevel, email string) *User {
	return &User{
		Level: level,
		Email: email,
	}
}

type UserSettings struct {
	PayloadReadTimeout uint32
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
