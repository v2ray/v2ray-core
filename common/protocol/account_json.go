// +build json

package protocol

import (
	"errors"

	"github.com/v2ray/v2ray-core/common/uuid"
)

type AccountJson struct {
	ID       string `json:"id"`
	AlterIds uint16 `json:"alterId"`

	Username string `json:"user"`
	Password string `json:"pass"`
}

func (this *AccountJson) GetAccount() (Account, error) {
	if len(this.ID) > 0 {
		id, err := uuid.ParseString(this.ID)
		if err != nil {
			return nil, err
		}

		primaryID := NewID(id)
		alterIDs := NewAlterIDs(primaryID, this.AlterIds)

		return &VMessAccount{
			ID:       primaryID,
			AlterIDs: alterIDs,
		}, nil
	}

	return nil, errors.New("Protocol: Malformed account.")
}
