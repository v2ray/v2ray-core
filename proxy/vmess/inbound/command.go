package inbound

import (
	"hash/fnv"

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

		tag := this.features.Detour.ToTag
		if this.space.HasInboundHandlerManager() {
			handlerManager := this.space.InboundHandlerManager()
			handler, availableMin := handlerManager.GetHandler(tag)
			inboundHandler, ok := handler.(*VMessInboundHandler)
			if ok {
				if availableMin > 255 {
					availableMin = 255
				}
				cmd = byte(1)
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

	if cmd == 0 || commandBytes.Len()+4 > 256 {
		buffer.AppendBytes(byte(0), byte(0))
	} else {
		buffer.AppendBytes(cmd, byte(commandBytes.Len()+4))
		fnv1hash := fnv.New32a()
		fnv1hash.Write(commandBytes.Value)
		hashValue := fnv1hash.Sum32()
		buffer.AppendBytes(byte(hashValue>>24), byte(hashValue>>16), byte(hashValue>>8), byte(hashValue))
		buffer.Append(commandBytes.Value)
	}
}
