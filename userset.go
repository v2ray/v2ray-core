package core

import (
	"encoding/base64"
)

type UserSet struct {
	validUserIds   []ID
	userIdsAskHash map[string]int
}

func NewUserSet() *UserSet {
	vuSet := new(UserSet)
	vuSet.validUserIds = make([]ID, 0, 16)
	vuSet.userIdsAskHash = make(map[string]int)
	return vuSet
}

func hashBytesToString(hash []byte) string {
	return base64.StdEncoding.EncodeToString(hash)
}

func (us *UserSet) AddUser(user User) error {
	id := user.Id
	us.validUserIds = append(us.validUserIds, id)

	idBase64 := hashBytesToString(id.Hash([]byte("ASK")))
	us.userIdsAskHash[idBase64] = len(us.validUserIds) - 1

	return nil
}

func (us UserSet) IsValidUserId(askHash []byte) (*ID, bool) {
	askBase64 := hashBytesToString(askHash)
	idIndex, found := us.userIdsAskHash[askBase64]
	if found {
		return &us.validUserIds[idIndex], true
	}
	return nil, false
}
