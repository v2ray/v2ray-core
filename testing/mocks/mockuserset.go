package mocks

import (
	"github.com/v2ray/v2ray-core"
)

type MockUserSet struct {
	UserIds    []core.ID
	UserHashes map[string]int
}

func (us *MockUserSet) AddUser(user core.User) error {
	us.UserIds = append(us.UserIds, user.Id)
	return nil
}

func (us *MockUserSet) GetUser(userhash []byte) (*core.ID, int64, bool) {
	idx, found := us.UserHashes[string(userhash)]
	if found {
		return &us.UserIds[idx], 1234, true
	}
	return nil, 0, false
}
