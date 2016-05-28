package protocol

import (
	"github.com/v2ray/v2ray-core/common/dice"
)

type Account interface {
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
