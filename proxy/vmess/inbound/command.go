package inbound

import (
	"github.com/v2ray/v2ray-core/common/alloc"
	"github.com/v2ray/v2ray-core/common/log"
	"github.com/v2ray/v2ray-core/common/serial"
	"github.com/v2ray/v2ray-core/proxy/vmess/command"
)

func (this *VMessInboundHandler) generateCommand(buffer *alloc.Buffer) {
	cmd := byte(0)
	commandBytes := alloc.NewSmallBuffer().Clear()
	defer commandBytes.Release()

	if this.features != nil && this.features.Detour != nil {
		cmd = byte(1)
		tag := this.features.Detour.ToTag
		if this.space.HasInboundHandlerManager() {
			handlerManager := this.space.InboundHandlerManager()
			handler, availableMin := handlerManager.GetHandler(tag)

			inboundHandler, ok := handler.(*VMessInboundHandler)
			if ok {
				if availableMin > 255 {
					availableMin = 255
				}
				log.Info("VMessIn: Pick detour handler for port ", inboundHandler.Port(), " for ", availableMin, " minutes.")
				user := inboundHandler.GetUser()
				saCmd := &command.SwitchAccount{
					Port:     inboundHandler.Port(),
					ID:       user.ID.UUID(),
					AlterIds: serial.Uint16Literal(len(user.AlterIDs)),
					Level:    user.Level,
					ValidMin: byte(availableMin),
				}
				saCmd.Marshal(commandBytes)
			}
		}
	}

	if commandBytes.Len() > 256 {
		buffer.AppendBytes(byte(0), byte(0))
	} else {
		buffer.AppendBytes(cmd, byte(commandBytes.Len()))
		buffer.Append(commandBytes.Value)
	}
}
