package vmess

import (
	"strings"
	"sync"
	"time"

	"v2ray.com/core/common"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/common/task"
)

const (
	updateInterval   = 10 * time.Second
	cacheDurationSec = 120
)

type user struct {
	user    protocol.MemoryUser
	lastSec protocol.Timestamp
}

type TimedUserValidator struct {
	sync.RWMutex
	users    []*user
	userHash map[[16]byte]indexTimePair
	hasher   protocol.IDHash
	baseTime protocol.Timestamp
	task     *task.Periodic
}

type indexTimePair struct {
	user    *user
	timeInc uint32
}

func NewTimedUserValidator(hasher protocol.IDHash) *TimedUserValidator {
	tuv := &TimedUserValidator{
		users:    make([]*user, 0, 16),
		userHash: make(map[[16]byte]indexTimePair, 1024),
		hasher:   hasher,
		baseTime: protocol.Timestamp(time.Now().Unix() - cacheDurationSec*2),
	}
	tuv.task = &task.Periodic{
		Interval: updateInterval,
		Execute: func() error {
			tuv.updateUserHash()
			return nil
		},
	}
	common.Must(tuv.task.Start())
	return tuv
}

func (v *TimedUserValidator) generateNewHashes(nowSec protocol.Timestamp, user *user) {
	var hashValue [16]byte
	genHashForID := func(id *protocol.ID) {
		idHash := v.hasher(id.Bytes())
		lastSec := user.lastSec
		if lastSec < nowSec-cacheDurationSec*2 {
			lastSec = nowSec - cacheDurationSec*2
		}
		for ts := lastSec; ts <= nowSec; ts++ {
			common.Must2(idHash.Write(ts.Bytes(nil)))
			idHash.Sum(hashValue[:0])
			idHash.Reset()

			v.userHash[hashValue] = indexTimePair{
				user:    user,
				timeInc: uint32(ts - v.baseTime),
			}
		}
	}

	account := user.user.Account.(*MemoryAccount)

	genHashForID(account.ID)
	for _, id := range account.AlterIDs {
		genHashForID(id)
	}
	user.lastSec = nowSec
}

func (v *TimedUserValidator) removeExpiredHashes(expire uint32) {
	for key, pair := range v.userHash {
		if pair.timeInc < expire {
			delete(v.userHash, key)
		}
	}
}

func (v *TimedUserValidator) updateUserHash() {
	now := time.Now()
	nowSec := protocol.Timestamp(now.Unix() + cacheDurationSec)
	v.Lock()
	defer v.Unlock()

	for _, user := range v.users {
		v.generateNewHashes(nowSec, user)
	}

	expire := protocol.Timestamp(now.Unix() - cacheDurationSec)
	if expire > v.baseTime {
		v.removeExpiredHashes(uint32(expire - v.baseTime))
	}
}

func (v *TimedUserValidator) Add(u *protocol.MemoryUser) error {
	v.Lock()
	defer v.Unlock()

	nowSec := time.Now().Unix()

	uu := &user{
		user:    *u,
		lastSec: protocol.Timestamp(nowSec - cacheDurationSec),
	}
	v.users = append(v.users, uu)
	v.generateNewHashes(protocol.Timestamp(nowSec+cacheDurationSec), uu)

	return nil
}

func (v *TimedUserValidator) Get(userHash []byte) (*protocol.MemoryUser, protocol.Timestamp, bool) {
	defer v.RUnlock()
	v.RLock()

	var fixedSizeHash [16]byte
	copy(fixedSizeHash[:], userHash)
	pair, found := v.userHash[fixedSizeHash]
	if found {
		var user protocol.MemoryUser
		user = pair.user.user
		return &user, protocol.Timestamp(pair.timeInc) + v.baseTime, true
	}
	return nil, 0, false
}

func (v *TimedUserValidator) Remove(email string) bool {
	v.Lock()
	defer v.Unlock()

	email = strings.ToLower(email)
	idx := -1
	for i, u := range v.users {
		if strings.ToLower(u.user.Email) == email {
			idx = i
			break
		}
	}
	if idx == -1 {
		return false
	}
	ulen := len(v.users)
	if idx < ulen {
		v.users[idx] = v.users[ulen-1]
		v.users[ulen-1] = nil
		v.users = v.users[:ulen-1]
	}
	return true
}

// Close implements common.Closable.
func (v *TimedUserValidator) Close() error {
	return v.task.Close()
}
