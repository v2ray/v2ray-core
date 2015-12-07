package testing

import (
	v2net "github.com/v2ray/v2ray-core/common/net"
)

type TestRule struct {
	Function func(v2net.Destination) bool
	TagValue string
}

func (this *TestRule) Apply(dest v2net.Destination) bool {
	return this.Function(dest)
}

func (this *TestRule) Tag() string {
	return this.TagValue
}
