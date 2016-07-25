package protocol

import (
	"github.com/v2ray/v2ray-core/common"
)

type UserValidator interface {
	common.Releasable

	Add(user *User) error
	Get(timeHash []byte) (*User, Timestamp, bool)
}
