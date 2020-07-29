package conf

import (
	"encoding/json"

	"github.com/golang/protobuf/proto"

	"v2ray.com/core/common/net"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/common/serial"
	"v2ray.com/core/proxy/vless"
	"v2ray.com/core/proxy/vless/inbound"
	"v2ray.com/core/proxy/vless/outbound"
)

type VLessInboundFallback struct {
	Addr *Address `json:"addr"`
	Port uint16   `json:"port"`
	Unix string   `json:"unix"`
}

type VLessInboundConfig struct {
	Users      []json.RawMessage     `json:"clients"`
	Decryption string                `json:"decryption"`
	Fallback   *VLessInboundFallback `json:"fallback"`
}

// Build implements Buildable
func (c *VLessInboundConfig) Build() (proto.Message, error) {

	config := new(inbound.Config)

	if c.Decryption != "none" {
		return nil, newError(`please add/set "decryption":"none" directly to every VLESS "settings"`)
	}
	config.Decryption = c.Decryption

	if c.Fallback != nil {
		if c.Fallback.Unix != "" {
			if c.Fallback.Unix[0] == '@' {
				c.Fallback.Unix = "\x00" + c.Fallback.Unix[1:]
			}
		} else {
			if c.Fallback.Port == 0 {
				return nil, newError(`please fill in a valid value for "port" in VLESS "fallback"`)
			}
		}
		if c.Fallback.Addr == nil {
			c.Fallback.Addr = &Address{
				Address: net.ParseAddress("127.0.0.1"),
			}
		}
		config.Fallback = &inbound.Fallback{
			Addr: c.Fallback.Addr.Build(),
			Port: uint32(c.Fallback.Port),
			Unix: c.Fallback.Unix,
		}
	}

	config.User = make([]*protocol.User, len(c.Users))
	for idx, rawData := range c.Users {
		user := new(protocol.User)
		if err := json.Unmarshal(rawData, user); err != nil {
			return nil, newError("invalid VLESS user").Base(err)
		}
		account := new(vless.Account)
		if err := json.Unmarshal(rawData, account); err != nil {
			return nil, newError("invalid VLESS user").Base(err)
		}

		if account.Schedulers != "" {
			return nil, newError(`VLESS attr "schedulers" is not available in this version`)
		}
		if account.Encryption != "" {
			return nil, newError(`VLESS attr "encryption" should not in inbound settings`)
		}

		user.Account = serial.ToTypedMessage(account)
		config.User[idx] = user
	}

	return config, nil
}

type VLessOutboundTarget struct {
	Address *Address          `json:"address"`
	Port    uint16            `json:"port"`
	Users   []json.RawMessage `json:"users"`
}

type VLessOutboundConfig struct {
	Receivers []*VLessOutboundTarget `json:"vnext"`
}

// Build implements Buildable
func (c *VLessOutboundConfig) Build() (proto.Message, error) {

	config := new(outbound.Config)

	if len(c.Receivers) == 0 {
		return nil, newError("0 VLESS receiver configured")
	}
	serverSpecs := make([]*protocol.ServerEndpoint, len(c.Receivers))
	for idx, rec := range c.Receivers {
		if len(rec.Users) == 0 {
			return nil, newError("0 user configured for VLESS outbound")
		}
		if rec.Address == nil {
			return nil, newError("address is not set in VLESS outbound config")
		}
		spec := &protocol.ServerEndpoint{
			Address: rec.Address.Build(),
			Port:    uint32(rec.Port),
		}
		for _, rawUser := range rec.Users {
			user := new(protocol.User)
			if err := json.Unmarshal(rawUser, user); err != nil {
				return nil, newError("invalid VLESS user").Base(err)
			}
			account := new(vless.Account)
			if err := json.Unmarshal(rawUser, account); err != nil {
				return nil, newError("invalid VLESS user").Base(err)
			}

			if account.Schedulers != "" {
				return nil, newError(`VLESS attr "schedulers" is not available in this version`)
			}
			if account.Encryption != "none" {
				return nil, newError(`please add/set "encryption":"none" for every VLESS user in "users"`)
			}

			user.Account = serial.ToTypedMessage(account)
			spec.User = append(spec.User, user)
		}
		serverSpecs[idx] = spec
	}
	config.Receiver = serverSpecs

	return config, nil
}
