package mocks

import (
	"github.com/v2ray/v2ray-core/proxy/vmess"
)

type StaticUser struct {
	id *vmess.ID
}

func (this *StaticUser) ID() *vmess.ID {
	return this.id
}

func (this *StaticUser) Level() vmess.UserLevel {
	return vmess.UserLevelUntrusted
}

type StaticUserSet struct {
}

func (us *StaticUserSet) AddUser(user vmess.User) error {
	return nil
}

func (us *StaticUserSet) GetUser(userhash []byte) (vmess.User, int64, bool) {
	id, _ := vmess.NewID("703e9102-eb57-499c-8b59-faf4f371bb21")
	return &StaticUser{id: id}, 0, true
}
