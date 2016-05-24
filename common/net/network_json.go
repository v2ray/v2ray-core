// +build json

package net

import (
	"encoding/json"

	"github.com/v2ray/v2ray-core/common/collect"
)

func (this *NetworkList) UnmarshalJSON(data []byte) error {
	var strlist collect.StringList
	if err := json.Unmarshal(data, &strlist); err != nil {
		return err
	}
	*this = NewNetworkList(strlist)
	return nil
}
