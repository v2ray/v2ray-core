// +build json

package dns

import (
	"encoding/json"

	"github.com/v2ray/v2ray-core/common/serial"
)

func (this *CacheConfig) UnmarshalJSON(data []byte) error {
	var strlist serial.StringLiteralList
	if err := json.Unmarshal(data, strlist); err != nil {
		return err
	}
	config := &CacheConfig{
		TrustedTags: make(map[serial.StringLiteral]bool, strlist.Len()),
	}
	for _, str := range strlist {
		config.TrustedTags[str.TrimSpace()] = true
	}
	return nil
}
