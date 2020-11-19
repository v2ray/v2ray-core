// +build !confonly

package vmess

import (
	"crypto/hmac"
	"crypto/sha256"
	"hash"
	"hash/crc64"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"v2ray.com/core/common"
	"v2ray.com/core/common/dice"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/common/serial"
	"v2ray.com/core/common/task"
	"v2ray.com/core/proxy/vmess/aead"
)

const (
	updateInterval   = 10 * time.Second
	cacheDurationSec = 120
)

type user struct {
	user    protocol.MemoryUser
	lastSec protocol.Timestamp
}

// TimedUserValidator is a user Validator based on time.
type TimedUserValidator struct {
	sync.RWMutex
	users    []*user
	userHash map[[16]byte]indexTimePair
	hasher   protocol.IDHash
	baseTime protocol.Timestamp
	task     *task.Periodic

	behaviorSeed  uint64
	behaviorFused bool

	aeadDecoderHolder *aead.AuthIDDecoderHolder
}

type indexTimePair struct {
	user    *user
	timeInc uint32

	taintedFuse *uint32
}

// NewTimedUserValidator creates a new TimedUserValidator.
func NewTimedUserValidator(hasher protocol.IDHash) *TimedUserValidator {
	tuv := &TimedUserValidator{
		users:             make([]*user, 0, 16),
		userHash:          make(map[[16]byte]indexTimePair, 1024),
		hasher:            hasher,
		baseTime:          protocol.Timestamp(time.Now().Unix() - cacheDurationSec*2),
		aeadDecoderHolder: aead.NewAuthIDDecoderHolder(),
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
	genEndSec := nowSec + cacheDurationSec
	genHashForID := func(id *protocol.ID) {
		idHash := v.hasher(id.Bytes())
		genBeginSec := user.lastSec
		if genBeginSec < nowSec-cacheDurationSec {
			genBeginSec = nowSec - cacheDurationSec
		}
		for ts := genBeginSec; ts <= genEndSec; ts++ {
			common.Must2(serial.WriteUint64(idHash, uint64(ts)))
			idHash.Sum(hashValue[:0])
			idHash.Reset()

			v.userHash[hashValue] = indexTimePair{
				user:        user,
				timeInc:     uint32(ts - v.baseTime),
				taintedFuse: new(uint32),
			}
		}
	}

	account := user.user.Account.(*MemoryAccount)

	genHashForID(account.ID)
	for _, id := range account.AlterIDs {
		genHashForID(id)
	}
	user.lastSec = genEndSec
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
	nowSec := protocol.Timestamp(now.Unix())
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
	v.generateNewHashes(protocol.Timestamp(nowSec), uu)

	account := uu.user.Account.(*MemoryAccount)
	if !v.behaviorFused {
		hashkdf := hmac.New(func() hash.Hash { return sha256.New() }, []byte("VMESSBSKDF"))
		hashkdf.Write(account.ID.Bytes())
		v.behaviorSeed = crc64.Update(v.behaviorSeed, crc64.MakeTable(crc64.ECMA), hashkdf.Sum(nil))
	}

	var cmdkeyfl [16]byte
	copy(cmdkeyfl[:], account.ID.CmdKey())
	v.aeadDecoderHolder.AddUser(cmdkeyfl, u)

	return nil
}

func (v *TimedUserValidator) Get(userHash []byte) (*protocol.MemoryUser, protocol.Timestamp, bool, error) {
	defer v.RUnlock()
	v.RLock()

	v.behaviorFused = true

	var fixedSizeHash [16]byte
	copy(fixedSizeHash[:], userHash)
	pair, found := v.userHash[fixedSizeHash]
	if found {
		user := pair.user.user
		if atomic.LoadUint32(pair.taintedFuse) == 0 {
			return &user, protocol.Timestamp(pair.timeInc) + v.baseTime, true, nil
		}
		return nil, 0, false, ErrTainted
	}
	return nil, 0, false, ErrNotFound
}

func (v *TimedUserValidator) GetAEAD(userHash []byte) (*protocol.MemoryUser, bool, error) {
	defer v.RUnlock()
	v.RLock()
	var userHashFL [16]byte
	copy(userHashFL[:], userHash)

	userd, err := v.aeadDecoderHolder.Match(userHashFL)
	if err != nil {
		return nil, false, err
	}
	return userd.(*protocol.MemoryUser), true, err
}

func (v *TimedUserValidator) Remove(email string) bool {
	v.Lock()
	defer v.Unlock()

	idx := -1
	for i := range v.users {
		if strings.EqualFold(v.users[i].user.Email, email) {
			idx = i
			var cmdkeyfl [16]byte
			copy(cmdkeyfl[:], v.users[i].user.Account.(*MemoryAccount).ID.CmdKey())
			v.aeadDecoderHolder.RemoveUser(cmdkeyfl)
			break
		}
	}
	if idx == -1 {
		return false
	}
	ulen := len(v.users)

	v.users[idx] = v.users[ulen-1]
	v.users[ulen-1] = nil
	v.users = v.users[:ulen-1]

	return true
}

// Close implements common.Closable.
func (v *TimedUserValidator) Close() error {
	return v.task.Close()
}

func (v *TimedUserValidator) GetBehaviorSeed() uint64 {
	v.Lock()
	defer v.Unlock()
	v.behaviorFused = true
	if v.behaviorSeed == 0 {
		v.behaviorSeed = dice.RollUint64()
	}
	return v.behaviorSeed
}

func (v *TimedUserValidator) BurnTaintFuse(userHash []byte) error {
	v.RLock()
	defer v.RUnlock()
	var userHashFL [16]byte
	copy(userHashFL[:], userHash)

	pair, found := v.userHash[userHashFL]
	if found {
		if atomic.CompareAndSwapUint32(pair.taintedFuse, 0, 1) {
			return nil
		}
		return ErrTainted
	}
	return ErrNotFound
}

var ErrNotFound = newError("Not Found")

var ErrTainted = newError("ErrTainted")
