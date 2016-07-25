package outbound

import (
	"time"

	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/common/protocol"
	"github.com/v2ray/v2ray-core/proxy/vmess"
)

func (this *VMessOutboundHandler) handleSwitchAccount(cmd *protocol.CommandSwitchAccount) {
	primary := protocol.NewID(cmd.ID)
	alters := protocol.NewAlterIDs(primary, cmd.AlterIds)
	account := &vmess.Account{
		ID:       primary,
		AlterIDs: alters,
	}
	user := protocol.NewUser(cmd.Level, "")
	user.Account = account
	dest := v2net.TCPDestination(cmd.Host, cmd.Port)
	until := time.Now().Add(time.Duration(cmd.ValidMin) * time.Minute)
	this.serverList.AddServer(protocol.NewServerSpec(dest, protocol.BeforeTime(until), user))
}

func (this *VMessOutboundHandler) handleCommand(dest v2net.Destination, cmd protocol.ResponseCommand) {
	switch typedCommand := cmd.(type) {
	case *protocol.CommandSwitchAccount:
		if typedCommand.Host == nil {
			typedCommand.Host = dest.Address()
		}
		this.handleSwitchAccount(typedCommand)
	default:
	}
}
