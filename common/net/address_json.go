// +build json

package net

import (
	"encoding/json"
	"net"
)

type AddressJson struct {
	Address Address
}

func (this *AddressJson) UnmarshalJSON(data []byte) error {
	var rawStr string
	if err := json.Unmarshal(data, &rawStr); err != nil {
		return err
	}
	ip := net.ParseIP(rawStr)
	if ip != nil {
		this.Address = IPAddress(ip)
	} else {
		this.Address = DomainAddress(rawStr)
	}
	return nil
}
