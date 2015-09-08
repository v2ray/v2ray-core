package core

import (
	"encoding/base64"
)

type VUserSet struct {
	validUserIds   []VID
	userIdsAskHash map[string]int
}

func NewVUserSet() *VUserSet {
	vuSet := new(VUserSet)
	vuSet.validUserIds = make([]VID, 0, 16)
	vuSet.userIdsAskHash = make(map[string]int)
	return vuSet
}

func hashBytesToString(hash []byte) string {
	return base64.StdEncoding.EncodeToString(hash)
}

func (us *VUserSet) AddUser(user VUser) error {
	id := user.Id
	us.validUserIds = append(us.validUserIds, id)

	idBase64 := hashBytesToString(id.Hash([]byte("ASK")))
	us.userIdsAskHash[idBase64] = len(us.validUserIds) - 1

	return nil
}

func (us VUserSet) IsValidUserId(askHash []byte) (*VID, bool) {
	askBase64 := hashBytesToString(askHash)
	idIndex, found := us.userIdsAskHash[askBase64]
	if found {
		return &us.validUserIds[idIndex], true
	}
	return nil, false
}
