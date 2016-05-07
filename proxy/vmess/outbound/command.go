package outbound

import (
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/common/protocol"
)

func (this *VMessOutboundHandler) handleSwitchAccount(cmd *protocol.CommandSwitchAccount) {
	primary := protocol.NewID(cmd.ID)
	alters := protocol.NewAlterIDs(primary, cmd.AlterIds.Value())
	user := protocol.NewUser(primary, alters, cmd.Level, "")
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
