// +build json

package net

import (
	"encoding/json"
)

func (this *AddressPB) UnmarshalJSON(data []byte) error {
	var rawStr string
	if err := json.Unmarshal(data, &rawStr); err != nil {
		return err
	}
	addr := ParseAddress(rawStr)
	switch addr.Family() {
	case AddressFamilyIPv4, AddressFamilyIPv6:
		this.Address = &AddressPB_Ip{
			Ip: []byte(addr.IP()),
		}
	case AddressFamilyDomain:
		this.Address = &AddressPB_Domain{
			Domain: addr.Domain(),
		}
	}

	return nil
}
