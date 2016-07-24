package protocol

import (
	"github.com/v2ray/v2ray-core/common/dice"
)

type Account interface {
	Equals(Account) bool
}

type VMessAccount struct {
	ID       *ID
	AlterIDs []*ID
}

func (this *VMessAccount) AnyValidID() *ID {
	if len(this.AlterIDs) == 0 {
		return this.ID
	}
	return this.AlterIDs[dice.Roll(len(this.AlterIDs))]
}

func (this *VMessAccount) Equals(account Account) bool {
	vmessAccount, ok := account.(*VMessAccount)
	if !ok {
		return false
	}
	// TODO: handle AlterIds difference
	return this.ID.Equals(vmessAccount.ID)
}
