package core

import (
	"time"
)

const (
	updateIntervalSec = 10
	cacheDurationSec  = 120
)

type UserSet interface {
	AddUser(user User) error
	GetUser(timeHash []byte) (*ID, bool)
}

type TimedUserSet struct {
	validUserIds []ID
	userHashes   map[string]int
}

type hashEntry struct {
	hash    string
	timeSec int64
}

func NewTimedUserSet() UserSet {
	vuSet := new(TimedUserSet)
	vuSet.validUserIds = make([]ID, 0, 16)
	vuSet.userHashes = make(map[string]int)

	go vuSet.updateUserHash(time.Tick(updateIntervalSec * time.Second))
	return vuSet
}

func (us *TimedUserSet) updateUserHash(tick <-chan time.Time) {
	now := time.Now().UTC()
	lastSec := now.Unix() - cacheDurationSec

	hash2Remove := make(chan hashEntry, updateIntervalSec*2)
	lastSec2Remove := now.Unix() + cacheDurationSec
	for {
		now := <-tick
		nowSec := now.UTC().Unix()

		remove2Sec := nowSec - cacheDurationSec
		if remove2Sec > lastSec2Remove {
			for lastSec2Remove+1 < remove2Sec {
				entry := <-hash2Remove
				lastSec2Remove = entry.timeSec
				delete(us.userHashes, entry.hash)
			}
		}

		for i := lastSec + 1; i <= nowSec; i++ {
			for idx, id := range us.validUserIds {
				idHash := id.TimeHash(i)
				hash2Remove <- hashEntry{string(idHash), i}
				us.userHashes[string(idHash)] = idx
			}
		}
	}
}

func (us *TimedUserSet) AddUser(user User) error {
	id := user.Id
	us.validUserIds = append(us.validUserIds, id)
	return nil
}

func (us TimedUserSet) GetUser(userHash []byte) (*ID, bool) {
	idIndex, found := us.userHashes[string(userHash)]
	if found {
		return &us.validUserIds[idIndex], true
	}
	return nil, false
}
