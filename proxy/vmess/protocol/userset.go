package protocol

import (
	"sync"
	"time"

	proto "github.com/v2ray/v2ray-core/common/protocol"
	"github.com/v2ray/v2ray-core/common/serial"
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
	id             *proto.ID
	userIdx        int
	lastSec        Timestamp
	lastSecRemoval Timestamp
}

type UserSet interface {
	AddUser(user *proto.User) error
	GetUser(timeHash []byte) (*proto.User, Timestamp, bool)
}

type TimedUserSet struct {
	validUsers []*proto.User
	userHash   map[[16]byte]*indexTimePair
	ids        []*idEntry
	access     sync.RWMutex
}

type indexTimePair struct {
	index   int
	timeSec Timestamp
}

func NewTimedUserSet() UserSet {
	tus := &TimedUserSet{
		validUsers: make([]*proto.User, 0, 16),
		userHash:   make(map[[16]byte]*indexTimePair, 512),
		access:     sync.RWMutex{},
		ids:        make([]*idEntry, 0, 512),
	}
	go tus.updateUserHash(time.Tick(updateIntervalSec * time.Second))
	return tus
}

func (us *TimedUserSet) generateNewHashes(nowSec Timestamp, idx int, entry *idEntry) {
	var hashValue [16]byte
	var hashValueRemoval [16]byte
	idHash := IDHash(entry.id.Bytes())
	for entry.lastSec <= nowSec {
		idHash.Write(entry.lastSec.Bytes())
		idHash.Sum(hashValue[:0])
		idHash.Reset()

		idHash.Write(entry.lastSecRemoval.Bytes())
		idHash.Sum(hashValueRemoval[:0])
		idHash.Reset()

		us.access.Lock()
		us.userHash[hashValue] = &indexTimePair{idx, entry.lastSec}
		delete(us.userHash, hashValueRemoval)
		us.access.Unlock()

		entry.lastSec++
		entry.lastSecRemoval++
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

func (us *TimedUserSet) AddUser(user *proto.User) error {
	idx := len(us.validUsers)
	us.validUsers = append(us.validUsers, user)

	nowSec := time.Now().Unix()

	entry := &idEntry{
		id:             user.ID,
		userIdx:        idx,
		lastSec:        Timestamp(nowSec - cacheDurationSec),
		lastSecRemoval: Timestamp(nowSec - cacheDurationSec*3),
	}
	us.generateNewHashes(Timestamp(nowSec+cacheDurationSec), idx, entry)
	us.ids = append(us.ids, entry)
	for _, alterid := range user.AlterIDs {
		entry := &idEntry{
			id:             alterid,
			userIdx:        idx,
			lastSec:        Timestamp(nowSec - cacheDurationSec),
			lastSecRemoval: Timestamp(nowSec - cacheDurationSec*3),
		}
		us.generateNewHashes(Timestamp(nowSec+cacheDurationSec), idx, entry)
		us.ids = append(us.ids, entry)
	}

	return nil
}

func (us *TimedUserSet) GetUser(userHash []byte) (*proto.User, Timestamp, bool) {
	defer us.access.RUnlock()
	us.access.RLock()
	var fixedSizeHash [16]byte
	copy(fixedSizeHash[:], userHash)
	pair, found := us.userHash[fixedSizeHash]
	if found {
		return us.validUsers[pair.index], pair.timeSec, true
	}
	return nil, 0, false
}
