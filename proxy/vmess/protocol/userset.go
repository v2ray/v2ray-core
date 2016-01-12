package protocol

import (
	"sync"
	"time"

	"github.com/v2ray/v2ray-core/common/collect"
	"github.com/v2ray/v2ray-core/common/serial"
	"github.com/v2ray/v2ray-core/proxy/vmess"
)

const (
	updateIntervalSec = 10
	cacheDurationSec  = 120
)

type Timestamp int64

func (this Timestamp) Bytes() []byte {
	return serial.Int64Literal(this).Bytes()
}

func (this Timestamp) HashBytes() []byte {
	once := this.Bytes()
	bytes := make([]byte, 0, 32)
	bytes = append(bytes, once...)
	bytes = append(bytes, once...)
	bytes = append(bytes, once...)
	bytes = append(bytes, once...)
	return bytes
}

type idEntry struct {
	id      *vmess.ID
	userIdx int
	lastSec Timestamp
	hashes  *collect.SizedQueue
}

type UserSet interface {
	AddUser(user vmess.User) error
	GetUser(timeHash []byte) (vmess.User, Timestamp, bool)
}

type TimedUserSet struct {
	validUsers []vmess.User
	userHash   map[string]indexTimePair
	ids        []*idEntry
	access     sync.RWMutex
}

type indexTimePair struct {
	index   int
	timeSec Timestamp
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

func (us *TimedUserSet) generateNewHashes(nowSec Timestamp, idx int, entry *idEntry) {
	for entry.lastSec <= nowSec {
		idHash := IDHash(entry.id.Bytes())
		idHash.Write(entry.lastSec.Bytes())
		idHashSlice := idHash.Sum(nil)
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
		nowSec := Timestamp(now.Unix() + cacheDurationSec)
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
		lastSec: Timestamp(nowSec - cacheDurationSec),
		hashes:  collect.NewSizedQueue(2*cacheDurationSec + 1),
	}
	us.generateNewHashes(Timestamp(nowSec+cacheDurationSec), idx, entry)
	us.ids = append(us.ids, entry)
	for _, alterid := range user.AlterIDs() {
		entry := &idEntry{
			id:      alterid,
			userIdx: idx,
			lastSec: Timestamp(nowSec - cacheDurationSec),
			hashes:  collect.NewSizedQueue(2*cacheDurationSec + 1),
		}
		us.generateNewHashes(Timestamp(nowSec+cacheDurationSec), idx, entry)
		us.ids = append(us.ids, entry)
	}

	return nil
}

func (us *TimedUserSet) GetUser(userHash []byte) (vmess.User, Timestamp, bool) {
	defer us.access.RUnlock()
	us.access.RLock()
	pair, found := us.userHash[string(userHash)]
	if found {
		return us.validUsers[pair.index], pair.timeSec, true
	}
	return nil, 0, false
}
