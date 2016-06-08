// +build json

package net

import (
	"encoding/json"
)

type AddressJson struct {
	Address Address
}

func (this *AddressJson) UnmarshalJSON(data []byte) error {
	var rawStr string
	if err := json.Unmarshal(data, &rawStr); err != nil {
		return err
	}
	this.Address = ParseAddress(rawStr)
	return nil
}
