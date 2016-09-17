package protocol

import (
	"sync"
	"time"

	"v2ray.com/core/common/dice"
	v2net "v2ray.com/core/common/net"
)

type ValidationStrategy interface {
	IsValid() bool
	Invalidate()
}

type AlwaysValidStrategy struct{}

func AlwaysValid() ValidationStrategy {
	return AlwaysValidStrategy{}
}

func (this AlwaysValidStrategy) IsValid() bool {
	return true
}

func (this AlwaysValidStrategy) Invalidate() {}

type TimeoutValidStrategy struct {
	until time.Time
}

func BeforeTime(t time.Time) ValidationStrategy {
	return &TimeoutValidStrategy{
		until: t,
	}
}

func (this *TimeoutValidStrategy) IsValid() bool {
	return this.until.After(time.Now())
}

func (this *TimeoutValidStrategy) Invalidate() {
	this.until = time.Time{}
}

type ServerSpec struct {
	sync.RWMutex
	dest       v2net.Destination
	users      []*User
	valid      ValidationStrategy
	newAccount NewAccountFactory
}

func NewServerSpec(newAccount NewAccountFactory, dest v2net.Destination, valid ValidationStrategy, users ...*User) *ServerSpec {
	return &ServerSpec{
		dest:       dest,
		users:      users,
		valid:      valid,
		newAccount: newAccount,
	}
}

func NewServerSpecFromPB(newAccount NewAccountFactory, spec ServerSpecPB) *ServerSpec {
	dest := v2net.TCPDestination(spec.Address.AsAddress(), v2net.Port(spec.Port))
	return NewServerSpec(newAccount, dest, AlwaysValid(), spec.Users...)
}

func (this *ServerSpec) Destination() v2net.Destination {
	return this.dest
}

func (this *ServerSpec) HasUser(user *User) bool {
	this.RLock()
	defer this.RUnlock()

	accountA, err := user.GetTypedAccount(this.newAccount())
	if err != nil {
		return false
	}
	for _, u := range this.users {
		accountB, err := u.GetTypedAccount(this.newAccount())
		if err == nil && accountA.Equals(accountB) {
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
	return this.valid.IsValid()
}

func (this *ServerSpec) Invalidate() {
	this.valid.Invalidate()
}
