package mocks

import (
	"github.com/v2ray/v2ray-core/proxy/vmess/config"
)

type StaticUser struct {
	id *config.ID
}

func (this *StaticUser) ID() *config.ID {
	return this.id
}

func (this *StaticUser) Level() config.UserLevel {
	return config.UserLevelUntrusted
}

type StaticUserSet struct {
}

func (us *StaticUserSet) AddUser(user config.User) error {
	return nil
}

func (us *StaticUserSet) GetUser(userhash []byte) (config.User, int64, bool) {
	id, _ := config.NewID("703e9102-eb57-499c-8b59-faf4f371bb21")
	return &StaticUser{id: id}, 0, true
}
