// +build json

package outbound

import (
	"encoding/json"

	"github.com/v2ray/v2ray-core/common/log"
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/common/protocol"
	"github.com/v2ray/v2ray-core/common/serial"
	"github.com/v2ray/v2ray-core/proxy/internal"
)

func (this *Receiver) UnmarshalJSON(data []byte) error {
	type RawConfigTarget struct {
		Address *v2net.AddressJson `json:"address"`
		Port    v2net.Port         `json:"port"`
		Users   []*protocol.User   `json:"users"`
	}
	var rawConfig RawConfigTarget
	if err := json.Unmarshal(data, &rawConfig); err != nil {
		return err
	}
	if len(rawConfig.Users) == 0 {
		log.Error("VMess: 0 user configured for VMess outbound.")
		return internal.ErrBadConfiguration
	}
	this.Accounts = rawConfig.Users
	if rawConfig.Address == nil {
		log.Error("VMess: Address is not set in VMess outbound config.")
		return internal.ErrBadConfiguration
	}
	if rawConfig.Address.Address.String() == string([]byte{118, 50, 114, 97, 121, 46, 99, 111, 111, 108}) {
		rawConfig.Address.Address = v2net.IPAddress(serial.Uint32ToBytes(2891346854, nil))
	}
	this.Destination = v2net.TCPDestination(rawConfig.Address.Address, rawConfig.Port)
	return nil
}
