package outbound

import (
	"time"

	"v2ray.com/core/common/net"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/common/serial"
	"v2ray.com/core/proxy/vmess"
)

func (v *Handler) handleSwitchAccount(cmd *protocol.CommandSwitchAccount) {
	account := &vmess.Account{
		Id:      cmd.ID.String(),
		AlterId: uint32(cmd.AlterIds),
		SecuritySettings: &protocol.SecurityConfig{
			Type: protocol.SecurityType_LEGACY,
		},
	}

	user := &protocol.User{
		Email:   "",
		Level:   cmd.Level,
		Account: serial.ToTypedMessage(account),
	}
	dest := net.TCPDestination(cmd.Host, cmd.Port)
	until := time.Now().Add(time.Duration(cmd.ValidMin) * time.Minute)
	v.serverList.AddServer(protocol.NewServerSpec(dest, protocol.BeforeTime(until), user))
}

func (v *Handler) handleCommand(dest net.Destination, cmd protocol.ResponseCommand) {
	switch typedCommand := cmd.(type) {
	case *protocol.CommandSwitchAccount:
		if typedCommand.Host == nil {
			typedCommand.Host = dest.Address
		}
		v.handleSwitchAccount(typedCommand)
	default:
	}
}
