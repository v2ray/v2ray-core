package mocks

import (
	"github.com/v2ray/v2ray-core/proxy/vmess/config"
)

type MockUserSet struct {
	UserIds    []*config.ID
	UserHashes map[string]int
	Timestamps map[string]int64
}

func (us *MockUserSet) AddUser(user config.User) error {
	us.UserIds = append(us.UserIds, user.ID())
	return nil
}

func (us *MockUserSet) GetUser(userhash []byte) (*config.ID, int64, bool) {
	idx, found := us.UserHashes[string(userhash)]
	if found {
		return us.UserIds[idx], us.Timestamps[string(userhash)], true
	}
	return nil, 0, false
}
