package mocks

import (
	"github.com/v2ray/v2ray-core/proxy/vmess/config"
)

type MockUserSet struct {
	Users      []config.User
	UserHashes map[string]int
	Timestamps map[string]int64
}

func (us *MockUserSet) AddUser(user config.User) error {
	us.Users = append(us.Users, user)
	return nil
}

func (us *MockUserSet) GetUser(userhash []byte) (config.User, int64, bool) {
	idx, found := us.UserHashes[string(userhash)]
	if found {
		return us.Users[idx], us.Timestamps[string(userhash)], true
	}
	return nil, 0, false
}
