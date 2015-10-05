package user

import (
	"sync"
	"time"

	"github.com/v2ray/v2ray-core/common/collect"
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
	validUserIds        []ID
	userHash            map[string]indexTimePair
	userHashDeleteQueue *collect.TimedQueue
	access              sync.RWMutex
}

type indexTimePair struct {
	index   int
	timeSec int64
}

func NewTimedUserSet() UserSet {
	tus := &TimedUserSet{
		validUserIds:        make([]ID, 0, 16),
		userHash:            make(map[string]indexTimePair, 512),
		userHashDeleteQueue: collect.NewTimedQueue(updateIntervalSec),
		access:              sync.RWMutex{},
	}
	go tus.updateUserHash(time.Tick(updateIntervalSec * time.Second))
	go tus.removeEntries(tus.userHashDeleteQueue.RemovedEntries())
	return tus
}

func (us *TimedUserSet) removeEntries(entries <-chan interface{}) {
	for {
		entry := <-entries
		us.access.Lock()
		delete(us.userHash, entry.(string))
		us.access.Unlock()
	}
}

func (us *TimedUserSet) generateNewHashes(lastSec, nowSec int64, idx int, id ID) {
	idHash := NewTimeHash(HMACHash{})
	for lastSec < nowSec+cacheDurationSec {
		idHash := idHash.Hash(id.Bytes[:], lastSec)
		us.access.Lock()
		us.userHash[string(idHash)] = indexTimePair{idx, lastSec}
		us.access.Unlock()
		us.userHashDeleteQueue.Add(string(idHash), lastSec+2*cacheDurationSec)
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
	defer us.access.RUnlock()
	us.access.RLock()
	pair, found := us.userHash[string(userHash)]
	if found {
		return &us.validUserIds[pair.index], pair.timeSec, true
	}
	return nil, 0, false
}
