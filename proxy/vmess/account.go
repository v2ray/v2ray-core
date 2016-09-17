package vmess

import (
	"v2ray.com/core/common/dice"
	"v2ray.com/core/common/log"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/common/uuid"
)

type Account struct {
	ID       *protocol.ID
	AlterIDs []*protocol.ID
}

func NewAccount() protocol.AsAccount {
	return &AccountPB{}
}

func (this *Account) AnyValidID() *protocol.ID {
	if len(this.AlterIDs) == 0 {
		return this.ID
	}
	return this.AlterIDs[dice.Roll(len(this.AlterIDs))]
}

func (this *Account) Equals(account protocol.Account) bool {
	vmessAccount, ok := account.(*Account)
	if !ok {
		return false
	}
	// TODO: handle AlterIds difference
	return this.ID.Equals(vmessAccount.ID)
}

func (this *AccountPB) AsAccount() (protocol.Account, error) {
	id, err := uuid.ParseString(this.Id)
	if err != nil {
		log.Error("VMess: Failed to parse ID: ", err)
		return nil, err
	}
	protoId := protocol.NewID(id)
	return &Account{
		ID:       protoId,
		AlterIDs: protocol.NewAlterIDs(protoId, uint16(this.AlterId)),
	}, nil
}
