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

type alwaysValidStrategy struct{}

func AlwaysValid() ValidationStrategy {
	return alwaysValidStrategy{}
}

func (alwaysValidStrategy) IsValid() bool {
	return true
}

func (alwaysValidStrategy) Invalidate() {}

type timeoutValidStrategy struct {
	until time.Time
}

func BeforeTime(t time.Time) ValidationStrategy {
	return &timeoutValidStrategy{
		until: t,
	}
}

func (s *timeoutValidStrategy) IsValid() bool {
	return s.until.After(time.Now())
}

func (s *timeoutValidStrategy) Invalidate() {
	s.until = time.Time{}
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

func (s *ServerSpec) Destination() net.Destination {
	return s.dest
}

func (s *ServerSpec) HasUser(user *User) bool {
	s.RLock()
	defer s.RUnlock()

	accountA, err := user.GetTypedAccount()
	if err != nil {
		return false
	}
	for _, u := range s.users {
		accountB, err := u.GetTypedAccount()
		if err == nil && accountA.Equals(accountB) {
			return true
		}
	}
	return false
}

func (s *ServerSpec) AddUser(user *User) {
	if s.HasUser(user) {
		return
	}

	s.Lock()
	defer s.Unlock()

	s.users = append(s.users, user)
}

func (s *ServerSpec) PickUser() *User {
	s.RLock()
	defer s.RUnlock()

	userCount := len(s.users)
	switch userCount {
	case 0:
		return nil
	case 1:
		return s.users[0]
	default:
		return s.users[dice.Roll(userCount)]
	}
}

func (s *ServerSpec) IsValid() bool {
	return s.valid.IsValid()
}

func (s *ServerSpec) Invalidate() {
	s.valid.Invalidate()
}
