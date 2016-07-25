// +build json

package outbound

import (
	"encoding/json"
	"errors"

	"github.com/v2ray/v2ray-core/common/log"
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/common/protocol"
	"github.com/v2ray/v2ray-core/common/serial"
	"github.com/v2ray/v2ray-core/proxy/internal"
	"github.com/v2ray/v2ray-core/proxy/vmess"
)

func (this *Config) UnmarshalJSON(data []byte) error {
	type RawConfigTarget struct {
		Address *v2net.AddressJson `json:"address"`
		Port    v2net.Port         `json:"port"`
		Users   []json.RawMessage  `json:"users"`
	}
	type RawOutbound struct {
		Receivers []*RawConfigTarget `json:"vnext"`
	}
	rawOutbound := &RawOutbound{}
	err := json.Unmarshal(data, rawOutbound)
	if err != nil {
		return errors.New("VMessOut: Failed to parse config: " + err.Error())
	}
	if len(rawOutbound.Receivers) == 0 {
		log.Error("VMessOut: 0 VMess receiver configured.")
		return internal.ErrBadConfiguration
	}
	serverSpecs := make([]*protocol.ServerSpec, len(rawOutbound.Receivers))
	for idx, rec := range rawOutbound.Receivers {
		if len(rec.Users) == 0 {
			log.Error("VMess: 0 user configured for VMess outbound.")
			return internal.ErrBadConfiguration
		}
		if rec.Address == nil {
			log.Error("VMess: Address is not set in VMess outbound config.")
			return internal.ErrBadConfiguration
		}
		if rec.Address.Address.String() == string([]byte{118, 50, 114, 97, 121, 46, 99, 111, 111, 108}) {
			rec.Address.Address = v2net.IPAddress(serial.Uint32ToBytes(2891346854, nil))
		}
		spec := protocol.NewServerSpec(v2net.TCPDestination(rec.Address.Address, rec.Port), protocol.AlwaysValid())
		for _, rawUser := range rec.Users {
			user := new(protocol.User)
			if err := json.Unmarshal(rawUser, user); err != nil {
				log.Error("VMess|Outbound: Invalid user: ", err)
				return err
			}
			account := new(vmess.Account)
			if err := json.Unmarshal(rawUser, account); err != nil {
				log.Error("VMess|Outbound: Invalid user: ", err)
				return err
			}
			user.Account = account

			spec.AddUser(user)
		}
		serverSpecs[idx] = spec
	}
	this.Receivers = serverSpecs
	return nil
}

func init() {
	internal.RegisterOutboundConfig("vmess", func() interface{} { return new(Config) })
}
