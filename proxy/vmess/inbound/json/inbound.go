package json

import (
	"github.com/v2ray/v2ray-core/proxy/internal/config"
	"github.com/v2ray/v2ray-core/proxy/internal/config/json"
	"github.com/v2ray/v2ray-core/proxy/vmess"
	vmessjson "github.com/v2ray/v2ray-core/proxy/vmess/json"
)

type Inbound struct {
	AllowedClients []*vmessjson.ConfigUser `json:"clients"`
}

func (c *Inbound) AllowedUsers() []vmess.User {
	users := make([]vmess.User, 0, len(c.AllowedClients))
	for _, rawUser := range c.AllowedClients {
		users = append(users, rawUser)
	}
	return users
}

func init() {
	config.RegisterInboundConnectionConfig("vmess", json.JsonConfigLoader(func() interface{} {
		return new(Inbound)
	}))
}
