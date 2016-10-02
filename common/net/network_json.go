// +build json

package net

import (
	"encoding/json"

	"v2ray.com/core/common/collect"
)

func (this *Network) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	*this = ParseNetwork(str)
	return nil
}

func (this *NetworkList) UnmarshalJSON(data []byte) error {
	var strlist collect.StringList
	if err := json.Unmarshal(data, &strlist); err != nil {
		return err
	}
	*this = *NewNetworkList(strlist)
	return nil
}
