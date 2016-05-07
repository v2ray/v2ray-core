package outbound

import (
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/common/protocol"
)

func (this *VMessOutboundHandler) handleSwitchAccount(cmd *protocol.CommandSwitchAccount) {
	user := protocol.NewUser(protocol.NewID(cmd.ID), cmd.Level, cmd.AlterIds.Value(), "")
	dest := v2net.TCPDestination(cmd.Host, cmd.Port)
	this.receiverManager.AddDetour(NewReceiver(dest, user), cmd.ValidMin)
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
