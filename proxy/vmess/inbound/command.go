package inbound

import (
	"github.com/v2ray/v2ray-core/common/log"
	"github.com/v2ray/v2ray-core/common/protocol"
	"github.com/v2ray/v2ray-core/proxy/vmess"
)

func (this *VMessInboundHandler) generateCommand(request *protocol.RequestHeader) protocol.ResponseCommand {
	if this.detours != nil {
		tag := this.detours.ToTag
		if this.inboundHandlerManager != nil {
			handler, availableMin := this.inboundHandlerManager.GetHandler(tag)
			inboundHandler, ok := handler.(*VMessInboundHandler)
			if ok {
				if availableMin > 255 {
					availableMin = 255
				}

				log.Info("VMessIn: Pick detour handler for port ", inboundHandler.Port(), " for ", availableMin, " minutes.")
				user := inboundHandler.GetUser(request.User.Email)
				if user == nil {
					return nil
				}
				return &protocol.CommandSwitchAccount{
					Port:     inboundHandler.Port(),
					ID:       user.Account.(*vmess.Account).ID.UUID(),
					AlterIds: uint16(len(user.Account.(*vmess.Account).AlterIDs)),
					Level:    user.Level,
					ValidMin: byte(availableMin),
				}
			}
		}
	}

	return nil
}
