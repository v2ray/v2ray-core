package protocol

import "time"

var (
	ErrAccountMissing     = newError("Account is not specified.")
	ErrNonMessageType     = newError("Not a protobuf message.")
	ErrUnknownAccountType = newError("Unknown account type.")
)

func (v *User) GetTypedAccount() (Account, error) {
	if v.GetAccount() == nil {
		return nil, ErrAccountMissing
	}

	rawAccount, err := v.Account.GetInstance()
	if err != nil {
		return nil, err
	}
	if asAccount, ok := rawAccount.(AsAccount); ok {
		return asAccount.AsAccount()
	}
	if account, ok := rawAccount.(Account); ok {
		return account, nil
	}
	return nil, newError("Unknown account type: ", v.Account.Type)
}

func (v *User) GetSettings() UserSettings {
	settings := UserSettings{}
	switch v.Level {
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
