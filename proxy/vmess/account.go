package vmess

import (
	"v2ray.com/core/common/dice"
	"v2ray.com/core/common/log"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/common/uuid"
)

type InternalAccount struct {
	ID       *protocol.ID
	AlterIDs []*protocol.ID
	Security protocol.Security
}

func (v *InternalAccount) AnyValidID() *protocol.ID {
	if len(v.AlterIDs) == 0 {
		return v.ID
	}
	return v.AlterIDs[dice.Roll(len(v.AlterIDs))]
}

func (v *InternalAccount) Equals(account protocol.Account) bool {
	vmessAccount, ok := account.(*InternalAccount)
	if !ok {
		return false
	}
	// TODO: handle AlterIds difference
	return v.ID.Equals(vmessAccount.ID)
}

func (v *Account) AsAccount() (protocol.Account, error) {
	id, err := uuid.ParseString(v.Id)
	if err != nil {
		log.Error("VMess: Failed to parse ID: ", err)
		return nil, err
	}
	protoId := protocol.NewID(id)
	return &InternalAccount{
		ID:       protoId,
		AlterIDs: protocol.NewAlterIDs(protoId, uint16(v.AlterId)),
		Security: v.SecuritySettings.AsSecurity(),
	}, nil
}
