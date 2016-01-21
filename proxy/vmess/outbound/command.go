package outbound

import (
	"github.com/v2ray/v2ray-core/common/log"
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/proxy/vmess"
	"github.com/v2ray/v2ray-core/proxy/vmess/command"
)

func (this *VMessOutboundHandler) handleSwitchAccount(cmd *command.SwitchAccount) {
	user := vmess.NewUser(vmess.NewID(cmd.ID), cmd.Level, cmd.AlterIds.Value())
	dest := v2net.TCPDestination(cmd.Host, cmd.Port)
	this.receiverManager.AddDetour(NewReceiver(dest, user), cmd.ValidMin)
}

func (this *VMessOutboundHandler) handleCommand(dest v2net.Destination, cmdId byte, data []byte) {
	cmd, err := command.CreateResponseCommand(cmdId)
	if err != nil {
		log.Warning("VMessOut: Unknown response command (", cmdId, "): ", err)
		return
	}
	if err := cmd.Unmarshal(data); err != nil {
		log.Warning("VMessOut: Failed to parse response command: ", err)
		return
	}
	switch typedCommand := cmd.(type) {
	case *command.SwitchAccount:
		if typedCommand.Host == nil {
			typedCommand.Host = dest.Address()
		}
		this.handleSwitchAccount(typedCommand)
	default:
	}
}
