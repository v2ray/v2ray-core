package internal

import (
	"github.com/v2ray/v2ray-core/common/serial"
)

type CacheConfig struct {
	TrustedTags map[serial.StringLiteral]bool
}

func (this *CacheConfig) IsTrustedSource(tag serial.StringLiteral) bool {
	_, found := this.TrustedTags[tag]
	return found
}
