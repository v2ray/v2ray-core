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

type idEntry struct {
	id      *vmess.ID
	userIdx int
	lastSec int64
	hashes  *collect.SizedQueue
}

type UserSet interface {
	AddUser(user vmess.User) error
	GetUser(timeHash []byte) (vmess.User, int64, bool)
}

type TimedUserSet struct {
	validUsers []vmess.User
	userHash   map[string]indexTimePair
	ids        []*idEntry
	access     sync.RWMutex
}

type indexTimePair struct {
	index   int
	timeSec int64
}

func NewTimedUserSet() UserSet {
	tus := &TimedUserSet{
		validUsers: make([]vmess.User, 0, 16),
		userHash:   make(map[string]indexTimePair, 512),
		access:     sync.RWMutex{},
		ids:        make([]*idEntry, 0, 512),
	}
	go tus.updateUserHash(time.Tick(updateIntervalSec * time.Second))
	return tus
}

func (us *TimedUserSet) generateNewHashes(nowSec int64, idx int, entry *idEntry) {
	idHash := NewTimeHash(HMACHash{})
	for entry.lastSec <= nowSec {
		idHashSlice := idHash.Hash(entry.id.Bytes(), entry.lastSec)
		hashValue := string(idHashSlice)
		us.access.Lock()
		us.userHash[hashValue] = indexTimePair{idx, entry.lastSec}
		us.access.Unlock()

		hash2Remove := entry.hashes.Put(hashValue)
		if hash2Remove != nil {
			us.access.Lock()
			delete(us.userHash, hash2Remove.(string))
			us.access.Unlock()
		}
		entry.lastSec++
	}
}

func (us *TimedUserSet) updateUserHash(tick <-chan time.Time) {
	for now := range tick {
		nowSec := now.Unix() + cacheDurationSec
		for _, entry := range us.ids {
			us.generateNewHashes(nowSec, entry.userIdx, entry)
		}
	}
}

func (us *TimedUserSet) AddUser(user vmess.User) error {
	idx := len(us.validUsers)
	us.validUsers = append(us.validUsers, user)

	nowSec := time.Now().Unix()

	entry := &idEntry{
		id:      user.ID(),
		userIdx: idx,
		lastSec: nowSec - cacheDurationSec,
		hashes:  collect.NewSizedQueue(2*cacheDurationSec + 1),
	}
	us.generateNewHashes(nowSec+cacheDurationSec, idx, entry)
	us.ids = append(us.ids, entry)
	for _, alterid := range user.AlterIDs() {
		entry := &idEntry{
			id:      alterid,
			userIdx: idx,
			lastSec: nowSec - cacheDurationSec,
			hashes:  collect.NewSizedQueue(2*cacheDurationSec + 1),
		}
		us.generateNewHashes(nowSec+cacheDurationSec, idx, entry)
		us.ids = append(us.ids, entry)
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
