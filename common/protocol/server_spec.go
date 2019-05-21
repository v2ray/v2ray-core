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
	users []*MemoryUser
	valid ValidationStrategy
}

func NewServerSpec(dest net.Destination, valid ValidationStrategy, users ...*MemoryUser) *ServerSpec {
	return &ServerSpec{
		dest:  dest,
		users: users,
		valid: valid,
	}
}

func NewServerSpecFromPB(spec ServerEndpoint) (*ServerSpec, error) {
	dest := net.TCPDestination(spec.Address.AsAddress(), net.Port(spec.Port))
	mUsers := make([]*MemoryUser, len(spec.User))
	for idx, u := range spec.User {
		mUser, err := u.ToMemoryUser()
		if err != nil {
			return nil, err
		}
		mUsers[idx] = mUser
	}
	return NewServerSpec(dest, AlwaysValid(), mUsers...), nil
}

func (s *ServerSpec) Destination() net.Destination {
	return s.dest
}

func (s *ServerSpec) HasUser(user *MemoryUser) bool {
	s.RLock()
	defer s.RUnlock()

	for _, u := range s.users {
		if u.Account.Equals(user.Account) {
			return true
		}
	}
	return false
}

func (s *ServerSpec) AddUser(user *MemoryUser) {
	if s.HasUser(user) {
		return
	}

	s.Lock()
	defer s.Unlock()

	s.users = append(s.users, user)
}

func (s *ServerSpec) PickUser() *MemoryUser {
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
