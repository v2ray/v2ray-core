package protocol

import (
	"v2ray.com/core/common"
)

type UserValidator interface {
	common.Releasable

	Add(user *User) error
	Get(timeHash []byte) (*User, Timestamp, bool)
}
