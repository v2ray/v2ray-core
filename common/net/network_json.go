// +build json

package net

import (
	"encoding/json"

	"github.com/v2ray/v2ray-core/common/serial"
)

func (this *NetworkList) UnmarshalJSON(data []byte) error {
	var strlist serial.StringLiteralList
	if err := json.Unmarshal(data, &strlist); err != nil {
		return err
	}
	*this = NewNetworkList(strlist)
	return nil
}
