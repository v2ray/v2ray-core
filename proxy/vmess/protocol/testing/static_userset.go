package mocks

import (
	"github.com/v2ray/v2ray-core/common/uuid"
	"github.com/v2ray/v2ray-core/proxy/vmess"
	"github.com/v2ray/v2ray-core/proxy/vmess/protocol"
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

func (this *StaticUser) AlterIDs() []*vmess.ID {
	return nil
}

func (this *StaticUser) AnyValidID() *vmess.ID {
	return this.id
}

type StaticUserSet struct {
}

func (us *StaticUserSet) AddUser(user vmess.User) error {
	return nil
}

func (us *StaticUserSet) GetUser(userhash []byte) (vmess.User, protocol.Timestamp, bool) {
	id, _ := uuid.ParseString("703e9102-eb57-499c-8b59-faf4f371bb21")
	return &StaticUser{id: vmess.NewID(id)}, 0, true
}
