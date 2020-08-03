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
	Xver uint16   `json:"xver"`
}

type VLessInboundConfig struct {
	Users       []json.RawMessage     `json:"clients"`
	Decryption  string                `json:"decryption"`
	Fallback    *VLessInboundFallback `json:"fallback"`
	Fallback_h2 *VLessInboundFallback `json:"fallback_h2"`
}

// Build implements Buildable
func (c *VLessInboundConfig) Build() (proto.Message, error) {

	config := new(inbound.Config)

	if c.Decryption != "none" {
		return nil, newError(`please add/set "decryption":"none" directly to every VLESS "settings"`)
	}
	config.Decryption = c.Decryption

	if c.Fallback != nil {
		if c.Fallback.Xver > 2 {
			return nil, newError(`VLESS "fallback": invalid PROXY protocol version, "xver" only accepts 0, 1, 2`)
		}
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
			Xver: uint32(c.Fallback.Xver),
		}
	}

	if c.Fallback_h2 != nil {
		if config.Fallback == nil {
			return nil, newError(`VLESS "fallback_h2" can't exist alone without "fallback"`)
		}
		if c.Fallback_h2.Xver > 2 {
			return nil, newError(`VLESS "fallback_h2": invalid PROXY protocol version, "xver" only accepts 0, 1, 2`)
		}
		if c.Fallback_h2.Unix != "" {
			if c.Fallback_h2.Unix[0] == '@' {
				c.Fallback_h2.Unix = "\x00" + c.Fallback_h2.Unix[1:]
			}
		} else {
			if c.Fallback_h2.Port == 0 {
				return nil, newError(`please fill in a valid value for "port" in VLESS "fallback_h2"`)
			}
		}
		if c.Fallback_h2.Addr == nil {
			c.Fallback_h2.Addr = &Address{
				Address: net.ParseAddress("127.0.0.1"),
			}
		}
		config.FallbackH2 = &inbound.FallbackH2{
			Addr: c.Fallback_h2.Addr.Build(),
			Port: uint32(c.Fallback_h2.Port),
			Unix: c.Fallback_h2.Unix,
			Xver: uint32(c.Fallback_h2.Xver),
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
