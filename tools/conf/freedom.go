package conf

import (
	"net"
	"strings"

	"v2ray.com/core/common/errors"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/common/serial"
	"v2ray.com/core/proxy/freedom"
)

type FreedomConfig struct {
	DomainStrategy string `json:"domainStrategy"`
	Timeout        uint32 `json:"timeout"`
	Redirect       string `json:"redirect"`
}

func (v *FreedomConfig) Build() (*serial.TypedMessage, error) {
	config := new(freedom.Config)
	config.DomainStrategy = freedom.Config_AS_IS
	domainStrategy := strings.ToLower(v.DomainStrategy)
	if domainStrategy == "useip" || domainStrategy == "use_ip" {
		config.DomainStrategy = freedom.Config_USE_IP
	}
	config.Timeout = v.Timeout
	if len(v.Redirect) > 0 {
		host, portStr, err := net.SplitHostPort(v.Redirect)
		if err != nil {
			return nil, errors.Base(err).Message("Config: Invalid redirect address: ", v.Redirect, ": ", err)
		}
		port, err := v2net.PortFromString(portStr)
		if err != nil {
			return nil, errors.Base(err).Message("Config: Invalid redirect port: ", v.Redirect, ": ", err)
		}
		config.DestinationOverride = &freedom.DestinationOverride{
			Server: &protocol.ServerEndpoint{
				Address: v2net.NewIPOrDomain(v2net.ParseAddress(host)),
				Port:    uint32(port),
			},
		}
	}
	return serial.ToTypedMessage(config), nil
}
