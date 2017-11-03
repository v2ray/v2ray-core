package protocol

import "time"

func (u *User) GetTypedAccount() (Account, error) {
	if u.GetAccount() == nil {
		return nil, newError("Account missing").AtWarning()
	}

	rawAccount, err := u.Account.GetInstance()
	if err != nil {
		return nil, err
	}
	if asAccount, ok := rawAccount.(AsAccount); ok {
		return asAccount.AsAccount()
	}
	if account, ok := rawAccount.(Account); ok {
		return account, nil
	}
	return nil, newError("Unknown account type: ", u.Account.Type)
}

func (u *User) GetSettings() UserSettings {
	settings := UserSettings{}
	switch u.Level {
	case 0:
		settings.PayloadTimeout = time.Second * 30
	case 1:
		settings.PayloadTimeout = time.Minute * 2
	default:
		settings.PayloadTimeout = time.Minute * 5
	}
	return settings
}

type UserSettings struct {
	PayloadTimeout time.Duration
}
