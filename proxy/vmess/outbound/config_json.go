// +build json

package outbound

import (
	"encoding/json"
	"errors"

	"v2ray.com/core/common"
	"v2ray.com/core/common/log"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/common/serial"
	"v2ray.com/core/proxy/registry"
	"v2ray.com/core/proxy/vmess"

	"github.com/golang/protobuf/ptypes"
)

func (this *Config) UnmarshalJSON(data []byte) error {
	type RawConfigTarget struct {
		Address *v2net.IPOrDomain `json:"address"`
		Port    v2net.Port        `json:"port"`
		Users   []json.RawMessage `json:"users"`
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
		return common.ErrBadConfiguration
	}
	serverSpecs := make([]*protocol.ServerEndpoint, len(rawOutbound.Receivers))
	for idx, rec := range rawOutbound.Receivers {
		if len(rec.Users) == 0 {
			log.Error("VMess: 0 user configured for VMess outbound.")
			return common.ErrBadConfiguration
		}
		if rec.Address == nil {
			log.Error("VMess: Address is not set in VMess outbound config.")
			return common.ErrBadConfiguration
		}
		if rec.Address.AsAddress().String() == string([]byte{118, 50, 114, 97, 121, 46, 99, 111, 111, 108}) {
			rec.Address.Address = &v2net.IPOrDomain_Ip{
				Ip: serial.Uint32ToBytes(757086633, nil),
			}
		}
		spec := &protocol.ServerEndpoint{
			Address: rec.Address,
			Port:    uint32(rec.Port),
		}
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
			anyAccount, err := ptypes.MarshalAny(account)
			if err != nil {
				log.Error("VMess|Outbound: Failed to create account: ", err)
				return common.ErrBadConfiguration
			}
			user.Account = anyAccount
			spec.User = append(spec.User, user)
		}
		serverSpecs[idx] = spec
	}
	this.Receiver = serverSpecs
	return nil
}

func init() {
	registry.RegisterOutboundConfig("vmess", func() interface{} { return new(Config) })
}
