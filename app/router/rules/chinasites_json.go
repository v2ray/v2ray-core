// +build json

package rules

import (
	"encoding/json"
	"v2ray.com/core/common/log"
)

func parseChinaSitesRule(data []byte) (*RoutingRule, error) {
	rawRule := new(JsonRule)
	err := json.Unmarshal(data, rawRule)
	if err != nil {
		log.Error("Router: Invalid router rule: ", err)
		return nil, err
	}
	return &RoutingRule{
		Tag:    rawRule.OutboundTag,
		Domain: chinaSitesDomains,
	}, nil
}
