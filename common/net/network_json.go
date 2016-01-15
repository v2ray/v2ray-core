// +build json

package net

import (
	serialjson "github.com/v2ray/v2ray-core/common/serial/json"
)

func (this *NetworkList) UnmarshalJSON(data []byte) error {
	strlist, err := serialjson.UnmarshalStringList(data)
	if err != nil {
		return err
	}
	*this = NewNetworkList(strlist)
	return nil
}
