package user

import (
	"time"

	"github.com/v2ray/v2ray-core/common/collect"
	"github.com/v2ray/v2ray-core/common/log"
)

const (
	updateIntervalSec = 10
	cacheDurationSec  = 120
)

type UserSet interface {
	AddUser(user User) error
	GetUser(timeHash []byte) (*ID, int64, bool)
}

type TimedUserSet struct {
	validUserIds []ID
	userHash     *collect.TimedStringMap
}

type indexTimePair struct {
	index   int
	timeSec int64
}

func NewTimedUserSet() UserSet {
	tus := &TimedUserSet{
		validUserIds: make([]ID, 0, 16),
		userHash:     collect.NewTimedStringMap(updateIntervalSec),
	}
	go tus.updateUserHash(time.Tick(updateIntervalSec * time.Second))
	return tus
}

func (us *TimedUserSet) generateNewHashes(lastSec, nowSec int64, idx int, id ID) {
	idHash := NewTimeHash(HMACHash{})
	for lastSec < nowSec+cacheDurationSec {
		idHash := idHash.Hash(id.Bytes[:], lastSec)
		log.Debug("Valid User Hash: %v", idHash)
		us.userHash.Set(string(idHash), indexTimePair{idx, lastSec}, lastSec+2*cacheDurationSec)
		lastSec++
	}
}

func (us *TimedUserSet) updateUserHash(tick <-chan time.Time) {
	now := time.Now().UTC()
	lastSec := now.Unix()

	for {
		now := <-tick
		nowSec := now.UTC().Unix()
		for idx, id := range us.validUserIds {
			us.generateNewHashes(lastSec, nowSec, idx, id)
		}
		lastSec = nowSec
	}
}

func (us *TimedUserSet) AddUser(user User) error {
	id := user.Id
	idx := len(us.validUserIds)
	us.validUserIds = append(us.validUserIds, id)

	nowSec := time.Now().UTC().Unix()
	lastSec := nowSec - cacheDurationSec
	us.generateNewHashes(lastSec, nowSec, idx, id)

	return nil
}

func (us TimedUserSet) GetUser(userHash []byte) (*ID, int64, bool) {
	rawPair, found := us.userHash.Get(string(userHash))
	if found {
		pair := rawPair.(indexTimePair)
		return &us.validUserIds[pair.index], pair.timeSec, true
	}
	return nil, 0, false
}
