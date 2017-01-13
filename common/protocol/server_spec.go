package protocol

import (
	"sync"
	"time"

	"v2ray.com/core/common/dice"
	"v2ray.com/core/common/net"
)

type ValidationStrategy interface {
	IsValid() bool
	Invalidate()
}

type AlwaysValidStrategy struct{}

func AlwaysValid() ValidationStrategy {
	return AlwaysValidStrategy{}
}

func (v AlwaysValidStrategy) IsValid() bool {
	return true
}

func (v AlwaysValidStrategy) Invalidate() {}

type TimeoutValidStrategy struct {
	until time.Time
}

func BeforeTime(t time.Time) ValidationStrategy {
	return &TimeoutValidStrategy{
		until: t,
	}
}

func (v *TimeoutValidStrategy) IsValid() bool {
	return v.until.After(time.Now())
}

func (v *TimeoutValidStrategy) Invalidate() {
	v.until = time.Time{}
}

type ServerSpec struct {
	sync.RWMutex
	dest  net.Destination
	users []*User
	valid ValidationStrategy
}

func NewServerSpec(dest net.Destination, valid ValidationStrategy, users ...*User) *ServerSpec {
	return &ServerSpec{
		dest:  dest,
		users: users,
		valid: valid,
	}
}

func NewServerSpecFromPB(spec ServerEndpoint) *ServerSpec {
	dest := net.TCPDestination(spec.Address.AsAddress(), net.Port(spec.Port))
	return NewServerSpec(dest, AlwaysValid(), spec.User...)
}

func (v *ServerSpec) Destination() net.Destination {
	return v.dest
}

func (v *ServerSpec) HasUser(user *User) bool {
	v.RLock()
	defer v.RUnlock()

	accountA, err := user.GetTypedAccount()
	if err != nil {
		return false
	}
	for _, u := range v.users {
		accountB, err := u.GetTypedAccount()
		if err == nil && accountA.Equals(accountB) {
			return true
		}
	}
	return false
}

func (v *ServerSpec) AddUser(user *User) {
	if v.HasUser(user) {
		return
	}

	v.Lock()
	defer v.Unlock()

	v.users = append(v.users, user)
}

func (v *ServerSpec) PickUser() *User {
	userCount := len(v.users)
	switch userCount {
	case 0:
		return nil
	case 1:
		return v.users[0]
	default:
		return v.users[dice.Roll(userCount)]
	}
}

func (v *ServerSpec) IsValid() bool {
	return v.valid.IsValid()
}

func (v *ServerSpec) Invalidate() {
	v.valid.Invalidate()
}
