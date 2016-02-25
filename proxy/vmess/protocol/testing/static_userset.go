package mocks

import (
	"github.com/v2ray/v2ray-core/common/protocol"
	"github.com/v2ray/v2ray-core/common/uuid"
)

type StaticUserSet struct {
}

func (us *StaticUserSet) Add(user *protocol.User) error {
	return nil
}

func (us *StaticUserSet) Get(userhash []byte) (*protocol.User, protocol.Timestamp, bool) {
	id, _ := uuid.ParseString("703e9102-eb57-499c-8b59-faf4f371bb21")
	return &protocol.User{
		ID: protocol.NewID(id),
	}, 0, true
}
