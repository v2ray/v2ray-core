// +build json

package rules

import (
	"encoding/json"

	"github.com/v2ray/v2ray-core/common/log"
)

func parseChinaIPRule(data []byte) (*Rule, error) {
	rawRule := new(JsonRule)
	err := json.Unmarshal(data, rawRule)
	if err != nil {
		log.Error("Router: Invalid router rule: ", err)
		return nil, err
	}
	return NewChinaIPRule(rawRule.OutboundTag), nil
}
