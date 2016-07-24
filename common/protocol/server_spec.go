package protocol

import (
	"sync"
	"time"

	"github.com/v2ray/v2ray-core/common/dice"
	v2net "github.com/v2ray/v2ray-core/common/net"
)

type ServerSpec struct {
	sync.RWMutex
	Destination v2net.Destination

	users []*User
}

func NewServerSpec(dest v2net.Destination, users ...*User) *ServerSpec {
	return &ServerSpec{
		Destination: dest,
		users:       users,
	}
}

func (this *ServerSpec) HasUser(user *User) bool {
	this.RLock()
	defer this.RUnlock()

	account := user.Account
	for _, u := range this.users {
		if u.Account.Equals(account) {
			return true
		}
	}
	return false
}

func (this *ServerSpec) AddUser(user *User) {
	if this.HasUser(user) {
		return
	}

	this.Lock()
	defer this.Unlock()

	this.users = append(this.users, user)
}

func (this *ServerSpec) PickUser() *User {
	userCount := len(this.users)
	return this.users[dice.Roll(userCount)]
}

func (this *ServerSpec) IsValid() bool {
	return true
}

func (this *ServerSpec) SetValid(b bool) {
}

type TimeoutServerSpec struct {
	*ServerSpec
	until time.Time
}

func NewTimeoutServerSpec(spec *ServerSpec, until time.Time) *TimeoutServerSpec {
	return &TimeoutServerSpec{
		ServerSpec: spec,
		until:      until,
	}
}

func (this *TimeoutServerSpec) IsValid() bool {
	return this.until.Before(time.Now())
}

func (this *TimeoutServerSpec) SetValid(b bool) {
	if !b {
		this.until = time.Time{}
	}
}
