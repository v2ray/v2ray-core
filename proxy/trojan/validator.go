package trojan

import (
	"sync"

	"v2ray.com/core/common/protocol"
)

// Validator stores valid trojan users
type Validator struct {
	users sync.Map
}

// Add a trojan user
func (v *Validator) Add(u *protocol.MemoryUser) error {
	user := u.Account.(*MemoryAccount)
	v.users.Store(hexString(user.Key), u)
	return nil
}

// Get user with hashed key, nil if user doesn't exist.
func (v *Validator) Get(hash string) *protocol.MemoryUser {
	u, _ := v.users.Load(hash)
	if u != nil {
		return u.(*protocol.MemoryUser)
	}
	return nil
}
