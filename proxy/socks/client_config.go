package socks

import (
	"v2ray.com/core/common/protocol"

	"github.com/golang/protobuf/ptypes"
	google_protobuf "github.com/golang/protobuf/ptypes/any"
)

func AccountEquals(a, b *google_protobuf.Any) bool {
	accountA := new(Account)
	if err := ptypes.UnmarshalAny(a, accountA); err != nil {
		return false
	}
	accountB := new(Account)
	if err := ptypes.UnmarshalAny(b, accountB); err != nil {
		return false
	}
	return accountA.Equals(accountB)
}

func (this *Account) AsAny() (*google_protobuf.Any, error) {
	return ptypes.MarshalAny(this)
}

type ClientConfig struct {
	Servers []*protocol.ServerSpec
}
