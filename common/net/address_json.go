// +build json

package net

import (
	"encoding/json"
)

func (this *IPOrDomain) UnmarshalJSON(data []byte) error {
	var rawStr string
	if err := json.Unmarshal(data, &rawStr); err != nil {
		return err
	}
	addr := ParseAddress(rawStr)
	switch addr.Family() {
	case AddressFamilyIPv4, AddressFamilyIPv6:
		this.Address = &IPOrDomain_Ip{
			Ip: []byte(addr.IP()),
		}
	case AddressFamilyDomain:
		this.Address = &IPOrDomain_Domain{
			Domain: addr.Domain(),
		}
	}

	return nil
}
