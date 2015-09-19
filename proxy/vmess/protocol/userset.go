package protocol

import (
	"container/heap"
	"time"

	"github.com/v2ray/v2ray-core/log"
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
	userHashes   map[string]indexTimePair
	hash2Remove  hashEntrySet
}

type indexTimePair struct {
	index   int
	timeSec int64
}

type hashEntry struct {
	hash    string
	timeSec int64
}

type hashEntrySet []*hashEntry

func (set hashEntrySet) Len() int {
	return len(set)
}

func (set hashEntrySet) Less(i, j int) bool {
	return set[i].timeSec < set[j].timeSec
}

func (set hashEntrySet) Swap(i, j int) {
	tmp := set[i]
	set[i] = set[j]
	set[j] = tmp
}

func (set *hashEntrySet) Push(value interface{}) {
	entry := value.(*hashEntry)
	*set = append(*set, entry)
}

func (set *hashEntrySet) Pop() interface{} {
	old := *set
	n := len(old)
	v := old[n-1]
	*set = old[:n-1]
	return v
}

func NewTimedUserSet() UserSet {
	vuSet := new(TimedUserSet)
	vuSet.validUserIds = make([]ID, 0, 16)
	vuSet.userHashes = make(map[string]indexTimePair)
	vuSet.hash2Remove = make(hashEntrySet, 0, cacheDurationSec*10)

	go vuSet.updateUserHash(time.Tick(updateIntervalSec * time.Second))
	return vuSet
}

func (us *TimedUserSet) generateNewHashes(lastSec, nowSec int64, idx int, id ID) {
	idHash := NewTimeHash(HMACHash{})
	for lastSec < nowSec+cacheDurationSec {

		idHash := idHash.Hash(id.Bytes, lastSec)
		log.Debug("Valid User Hash: %v", idHash)
		heap.Push(&us.hash2Remove, &hashEntry{string(idHash), lastSec})
		us.userHashes[string(idHash)] = indexTimePair{idx, lastSec}
		lastSec++
	}
}

func (us *TimedUserSet) updateUserHash(tick <-chan time.Time) {
	now := time.Now().UTC()
	lastSec := now.Unix()
	lastSec2Remove := now.Unix()

	for {
		now := <-tick
		nowSec := now.UTC().Unix()

		remove2Sec := nowSec - cacheDurationSec
		if remove2Sec > lastSec2Remove {
			for lastSec2Remove+1 < remove2Sec {
				front := heap.Pop(&us.hash2Remove)
				entry := front.(*hashEntry)
				lastSec2Remove = entry.timeSec
				delete(us.userHashes, entry.hash)
			}
		}
		for idx, id := range us.validUserIds {
			us.generateNewHashes(lastSec, nowSec, idx, id)
		}
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
	pair, found := us.userHashes[string(userHash)]
	if found {
		return &us.validUserIds[pair.index], pair.timeSec, true
	}
	return nil, 0, false
}
