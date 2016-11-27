package inbound

import (
	"v2ray.com/core/common/log"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/proxy/vmess"
)

func (v *VMessInboundHandler) generateCommand(request *protocol.RequestHeader) protocol.ResponseCommand {
	if v.detours != nil {
		tag := v.detours.To
		if v.inboundHandlerManager != nil {
			handler, availableMin := v.inboundHandlerManager.GetHandler(tag)
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
				account, _ := user.GetTypedAccount()
				return &protocol.CommandSwitchAccount{
					Port:     inboundHandler.Port(),
					ID:       account.(*vmess.InternalAccount).ID.UUID(),
					AlterIds: uint16(len(account.(*vmess.InternalAccount).AlterIDs)),
					Level:    user.Level,
					ValidMin: byte(availableMin),
				}
			}
		}
	}

	return nil
}
