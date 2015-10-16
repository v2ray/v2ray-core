package json

import (
	"github.com/v2ray/v2ray-core/config"
	"github.com/v2ray/v2ray-core/config/json"
	vmessconfig "github.com/v2ray/v2ray-core/proxy/vmess/config"
)

type Inbound struct {
	AllowedClients []*ConfigUser `json:"clients"`
	UDP            bool          `json:"udp"`
}

func (c *Inbound) AllowedUsers() []vmessconfig.User {
	users := make([]vmessconfig.User, 0, len(c.AllowedClients))
	for _, rawUser := range c.AllowedClients {
		users = append(users, rawUser)
	}
	return users
}

func (c *Inbound) UDPEnabled() bool {
	return c.UDP
}

func init() {
	json.RegisterConfigType("vmess", config.TypeInbound, func() interface{} {
		return new(Inbound)
	})
}
