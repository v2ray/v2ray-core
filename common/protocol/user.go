package protocol

import (
	"errors"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
)

var (
	ErrUserMissing    = errors.New("User is not specified.")
	ErrAccountMissing = errors.New("Account is not specified.")
	ErrNonMessageType = errors.New("Not a protobuf message.")
)

func (this *User) GetTypedAccount(account AsAccount) (Account, error) {
	anyAccount := this.GetAccount()
	if anyAccount == nil {
		return nil, ErrAccountMissing
	}
	protoAccount, ok := account.(proto.Message)
	if !ok {
		return nil, ErrNonMessageType
	}
	err := ptypes.UnmarshalAny(anyAccount, protoAccount)
	if err != nil {
		return nil, err
	}
	return account.AsAccount()
}

func (this *User) GetSettings() UserSettings {
	settings := UserSettings{
		PayloadReadTimeout: 120,
	}
	if this.Level > 0 {
		settings.PayloadReadTimeout = 0
	}
	return settings
}

type UserSettings struct {
	PayloadReadTimeout uint32
}
