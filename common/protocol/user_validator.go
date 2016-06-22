package protocol

import (
	"sync"
	"time"

	"github.com/v2ray/v2ray-core/common"
	"github.com/v2ray/v2ray-core/common/signal"
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
	common.Releasable

	Add(user *User) error
	Get(timeHash []byte) (*User, Timestamp, bool)
}

type TimedUserValidator struct {
	sync.RWMutex
	running    bool
	validUsers []*User
	userHash   map[[16]byte]*indexTimePair
	ids        []*idEntry
	hasher     IDHash
	cancel     *signal.CancelSignal
}

type indexTimePair struct {
	index   int
	timeSec Timestamp
}

func NewTimedUserValidator(hasher IDHash) UserValidator {
	tus := &TimedUserValidator{
		validUsers: make([]*User, 0, 16),
		userHash:   make(map[[16]byte]*indexTimePair, 512),
		ids:        make([]*idEntry, 0, 512),
		hasher:     hasher,
		running:    true,
		cancel:     signal.NewCloseSignal(),
	}
	go tus.updateUserHash(updateIntervalSec * time.Second)
	return tus
}

func (this *TimedUserValidator) Release() {
	if !this.running {
		return
	}

	this.cancel.Cancel()
	<-this.cancel.WaitForDone()

	this.Lock()
	defer this.Unlock()

	if !this.running {
		return
	}

	this.running = false
	this.validUsers = nil
	this.userHash = nil
	this.ids = nil
	this.hasher = nil
	this.cancel = nil
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

		this.Lock()
		this.userHash[hashValue] = &indexTimePair{idx, entry.lastSec}
		delete(this.userHash, hashValueRemoval)
		this.Unlock()

		entry.lastSec++
		entry.lastSecRemoval++
	}
}

func (this *TimedUserValidator) updateUserHash(interval time.Duration) {
L:
	for {
		select {
		case now := <-time.After(interval):
			nowSec := Timestamp(now.Unix() + cacheDurationSec)
			for _, entry := range this.ids {
				this.generateNewHashes(nowSec, entry.userIdx, entry)
			}
		case <-this.cancel.WaitForCancel():
			break L
		}
	}
	this.cancel.Done()
}

func (this *TimedUserValidator) Add(user *User) error {
	idx := len(this.validUsers)
	this.validUsers = append(this.validUsers, user)
	account := user.Account.(*VMessAccount)

	nowSec := time.Now().Unix()

	entry := &idEntry{
		id:             account.ID,
		userIdx:        idx,
		lastSec:        Timestamp(nowSec - cacheDurationSec),
		lastSecRemoval: Timestamp(nowSec - cacheDurationSec*3),
	}
	this.generateNewHashes(Timestamp(nowSec+cacheDurationSec), idx, entry)
	this.ids = append(this.ids, entry)
	for _, alterid := range account.AlterIDs {
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
	defer this.RUnlock()
	this.RLock()

	if !this.running {
		return nil, 0, false
	}
	var fixedSizeHash [16]byte
	copy(fixedSizeHash[:], userHash)
	pair, found := this.userHash[fixedSizeHash]
	if found {
		return this.validUsers[pair.index], pair.timeSec, true
	}
	return nil, 0, false
}
