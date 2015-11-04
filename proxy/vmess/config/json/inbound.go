package json

import (
	"github.com/v2ray/v2ray-core/proxy/common/config/json"
	vmessconfig "github.com/v2ray/v2ray-core/proxy/vmess/config"
)

type Inbound struct {
	AllowedClients []*ConfigUser `json:"clients"`
}

func (c *Inbound) AllowedUsers() []vmessconfig.User {
	users := make([]vmessconfig.User, 0, len(c.AllowedClients))
	for _, rawUser := range c.AllowedClients {
		users = append(users, rawUser)
	}
	return users
}

func init() {
	json.RegisterInboundConnectionConfig("vmess", func() interface{} {
		return new(Inbound)
	})
}
