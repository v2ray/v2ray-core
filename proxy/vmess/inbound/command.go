package inbound

import (
	"github.com/v2ray/v2ray-core/common/alloc"
	"github.com/v2ray/v2ray-core/common/serial"
	"github.com/v2ray/v2ray-core/proxy/vmess/command"
)

func (this *VMessInboundHandler) generateCommand(buffer *alloc.Buffer) {
	cmd := byte(0)
	commandBytes := alloc.NewSmallBuffer().Clear()
	defer commandBytes.Release()

	if this.features != nil && this.features.Detour != nil {
		tag := this.features.Detour.ToTag
		if this.space.HasInboundHandlerManager() {
			handlerManager := this.space.InboundHandlerManager()
			handler, availableSec := handlerManager.GetHandler(tag)
			inboundHandler, ok := handler.(*VMessInboundHandler)
			if ok {
				user := inboundHandler.GetUser()
				availableSecUint16 := uint16(65535)
				if availableSec < 65535 {
					availableSecUint16 = uint16(availableSec)
				}

				saCmd := &command.SwitchAccount{
					Port:     inboundHandler.Port(),
					ID:       user.ID.UUID(),
					AlterIds: serial.Uint16Literal(len(user.AlterIDs)),
					ValidSec: serial.Uint16Literal(availableSecUint16),
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
