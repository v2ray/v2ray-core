// +build !confonly

package vless

import (
	"strings"
	"sync"

	"v2ray.com/core/common/protocol"
	"v2ray.com/core/common/uuid"
)

type Validator struct {
	// Considering email's usage here, map + sync.Mutex/RWMutex may have better performance.
	email sync.Map
	users sync.Map
}

func (v *Validator) Add(u *protocol.MemoryUser) error {
	if u.Email != "" {
		_, loaded := v.email.LoadOrStore(strings.ToLower(u.Email), u)
		if loaded {
			return newError("User ", u.Email, " already exists.")
		}
	}
	v.users.Store(u.Account.(*MemoryAccount).ID.UUID(), u)
	return nil
}

func (v *Validator) Del(e string) error {
	if e == "" {
		return newError("Email must not be empty.")
	}
	le := strings.ToLower(e)
	u, _ := v.email.Load(le)
	if u == nil {
		return newError("User ", e, " not found.")
	}
	v.email.Delete(le)
	v.users.Delete(u.(*protocol.MemoryUser).Account.(*MemoryAccount).ID.UUID())
	return nil
}

func (v *Validator) Get(id uuid.UUID) *protocol.MemoryUser {
	u, _ := v.users.Load(id)
	if u != nil {
		return u.(*protocol.MemoryUser)
	}
	return nil
}
