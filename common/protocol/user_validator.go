package protocol

import (
	"sync"
	"time"
)

const (
	updateIntervalSec = 10
	cacheDurationSec  = 120
)

type idEntry struct {
	id             *ID
	userIdx        int
	lastSec        Timestamp
	lastSecRemoval Timestamp
}

type UserValidator interface {
	Add(user *User) error
	Get(timeHash []byte) (*User, Timestamp, bool)
}

type TimedUserValidator struct {
	validUsers []*User
	userHash   map[[16]byte]*indexTimePair
	ids        []*idEntry
	access     sync.RWMutex
	hasher     IDHash
}

type indexTimePair struct {
	index   int
	timeSec Timestamp
}

func NewTimedUserValidator(hasher IDHash) UserValidator {
	tus := &TimedUserValidator{
		validUsers: make([]*User, 0, 16),
		userHash:   make(map[[16]byte]*indexTimePair, 512),
		access:     sync.RWMutex{},
		ids:        make([]*idEntry, 0, 512),
		hasher:     hasher,
	}
	go tus.updateUserHash(time.Tick(updateIntervalSec * time.Second))
	return tus
}

func (this *TimedUserValidator) generateNewHashes(nowSec Timestamp, idx int, entry *idEntry) {
	var hashValue [16]byte
	var hashValueRemoval [16]byte
	idHash := this.hasher(entry.id.Bytes())
	for entry.lastSec <= nowSec {
		idHash.Write(entry.lastSec.Bytes())
		idHash.Sum(hashValue[:0])
		idHash.Reset()

		idHash.Write(entry.lastSecRemoval.Bytes())
		idHash.Sum(hashValueRemoval[:0])
		idHash.Reset()

		this.access.Lock()
		this.userHash[hashValue] = &indexTimePair{idx, entry.lastSec}
		delete(this.userHash, hashValueRemoval)
		this.access.Unlock()

		entry.lastSec++
		entry.lastSecRemoval++
	}
}

func (this *TimedUserValidator) updateUserHash(tick <-chan time.Time) {
	for now := range tick {
		nowSec := Timestamp(now.Unix() + cacheDurationSec)
		for _, entry := range this.ids {
			this.generateNewHashes(nowSec, entry.userIdx, entry)
		}
	}
}

func (this *TimedUserValidator) Add(user *User) error {
	idx := len(this.validUsers)
	this.validUsers = append(this.validUsers, user)

	nowSec := time.Now().Unix()

	entry := &idEntry{
		id:             user.ID,
		userIdx:        idx,
		lastSec:        Timestamp(nowSec - cacheDurationSec),
		lastSecRemoval: Timestamp(nowSec - cacheDurationSec*3),
	}
	this.generateNewHashes(Timestamp(nowSec+cacheDurationSec), idx, entry)
	this.ids = append(this.ids, entry)
	for _, alterid := range user.AlterIDs {
		entry := &idEntry{
			id:             alterid,
			userIdx:        idx,
			lastSec:        Timestamp(nowSec - cacheDurationSec),
			lastSecRemoval: Timestamp(nowSec - cacheDurationSec*3),
		}
		this.generateNewHashes(Timestamp(nowSec+cacheDurationSec), idx, entry)
		this.ids = append(this.ids, entry)
	}

	return nil
}

func (this *TimedUserValidator) Get(userHash []byte) (*User, Timestamp, bool) {
	defer this.access.RUnlock()
	this.access.RLock()
	var fixedSizeHash [16]byte
	copy(fixedSizeHash[:], userHash)
	pair, found := this.userHash[fixedSizeHash]
	if found {
		return this.validUsers[pair.index], pair.timeSec, true
	}
	return nil, 0, false
}
