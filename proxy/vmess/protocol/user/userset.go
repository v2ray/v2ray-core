package user

import (
	"sync"
	"time"

	"github.com/v2ray/v2ray-core/common/collect"
	"github.com/v2ray/v2ray-core/proxy/vmess"
)

const (
	updateIntervalSec = 10
	cacheDurationSec  = 120
)

type UserSet interface {
	AddUser(user vmess.User) error
	GetUser(timeHash []byte) (vmess.User, int64, bool)
}

type TimedUserSet struct {
	validUsers          []vmess.User
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
		validUsers:          make([]vmess.User, 0, 16),
		userHash:            make(map[string]indexTimePair, 512),
		userHashDeleteQueue: collect.NewTimedQueue(updateIntervalSec),
		access:              sync.RWMutex{},
	}
	go tus.updateUserHash(time.Tick(updateIntervalSec * time.Second))
	go tus.removeEntries(tus.userHashDeleteQueue.RemovedEntries())
	return tus
}

func (us *TimedUserSet) removeEntries(entries <-chan interface{}) {
	for entry := range entries {
		us.access.Lock()
		delete(us.userHash, entry.(string))
		us.access.Unlock()
	}
}

func (us *TimedUserSet) generateNewHashes(lastSec, nowSec int64, idx int, id *vmess.ID) {
	idHash := NewTimeHash(HMACHash{})
	for lastSec < nowSec {
		idHash := idHash.Hash(id.Bytes(), lastSec)
		us.access.Lock()
		us.userHash[string(idHash)] = indexTimePair{idx, lastSec}
		us.access.Unlock()
		us.userHashDeleteQueue.Add(string(idHash), lastSec+2*cacheDurationSec)
		lastSec++
	}
}

func (us *TimedUserSet) updateUserHash(tick <-chan time.Time) {
	lastSec := time.Now().Unix() - cacheDurationSec

	for now := range tick {
		nowSec := now.Unix() + cacheDurationSec
		for idx, user := range us.validUsers {
			us.generateNewHashes(lastSec, nowSec, idx, user.ID())
			for _, alterid := range user.AlterIDs() {
				us.generateNewHashes(lastSec, nowSec, idx, alterid)
			}
		}
		lastSec = nowSec
	}
}

func (us *TimedUserSet) AddUser(user vmess.User) error {
	id := user.ID()
	idx := len(us.validUsers)
	us.validUsers = append(us.validUsers, user)

	nowSec := time.Now().Unix()
	lastSec := nowSec - cacheDurationSec
	us.generateNewHashes(lastSec, nowSec+cacheDurationSec, idx, id)
	for _, alterid := range user.AlterIDs() {
		us.generateNewHashes(lastSec, nowSec+cacheDurationSec, idx, alterid)
	}

	return nil
}

func (us *TimedUserSet) GetUser(userHash []byte) (vmess.User, int64, bool) {
	defer us.access.RUnlock()
	us.access.RLock()
	pair, found := us.userHash[string(userHash)]
	if found {
		return us.validUsers[pair.index], pair.timeSec, true
	}
	return nil, 0, false
}
