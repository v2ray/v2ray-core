package mocks

import (
	proto "github.com/v2ray/v2ray-core/common/protocol"
	"github.com/v2ray/v2ray-core/common/uuid"
	"github.com/v2ray/v2ray-core/proxy/vmess/protocol"
)

type StaticUserSet struct {
}

func (us *StaticUserSet) AddUser(user *proto.User) error {
	return nil
}

func (us *StaticUserSet) GetUser(userhash []byte) (*proto.User, protocol.Timestamp, bool) {
	id, _ := uuid.ParseString("703e9102-eb57-499c-8b59-faf4f371bb21")
	return &proto.User{
		ID: proto.NewID(id),
	}, 0, true
}
