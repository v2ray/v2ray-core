package outbound

import (
	"math/rand"
	"sync"
	"time"

	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/proxy/vmess"
)

type Receiver struct {
	sync.RWMutex
	Destination v2net.Destination
	Accounts    []*vmess.User
}

func NewReceiver(dest v2net.Destination, users ...*vmess.User) *Receiver {
	return &Receiver{
		Destination: dest,
		Accounts:    users,
	}
}

func (this *Receiver) HasUser(user *vmess.User) bool {
	this.RLock()
	defer this.RUnlock()
	for _, u := range this.Accounts {
		// TODO: handle AlterIds difference.
		if u.ID.Equals(user.ID) {
			return true
		}
	}
	return false
}

func (this *Receiver) AddUser(user *vmess.User) {
	if this.HasUser(user) {
		return
	}
	this.Lock()
	this.Accounts = append(this.Accounts, user)
	this.Unlock()
}

func (this *Receiver) PickUser() *vmess.User {
	userLen := len(this.Accounts)
	userIdx := 0
	if userLen > 1 {
		userIdx = rand.Intn(userLen)
	}
	return this.Accounts[userIdx]
}

type ExpiringReceiver struct {
	*Receiver
	until time.Time
}

func (this *ExpiringReceiver) Expired() bool {
	return this.until.After(time.Now())
}

type ReceiverManager struct {
	receivers    []*Receiver
	detours      []*ExpiringReceiver
	detourAccess sync.RWMutex
}

func NewReceiverManager(receivers []*Receiver) *ReceiverManager {
	return &ReceiverManager{
		receivers: receivers,
		detours:   make([]*ExpiringReceiver, 0, 16),
	}
}

func (this *ReceiverManager) AddDetour(rec *Receiver, availableMin byte) {
	if availableMin < 2 {
		return
	}
	this.detourAccess.RLock()
	destExists := false
	for _, r := range this.detours {
		if r.Destination == rec.Destination {
			destExists = true
			// Destination exists, add new user if necessary
			for _, u := range rec.Accounts {
				r.AddUser(u)
			}
		}
	}

	this.detourAccess.RUnlock()
	expRec := &ExpiringReceiver{
		Receiver: rec,
		until:    time.Now().Add(time.Duration(availableMin-1) * time.Minute),
	}
	if !destExists {
		this.detourAccess.Lock()
		this.detours = append(this.detours, expRec)
		this.detourAccess.Unlock()
	}
}

func (this *ReceiverManager) pickDetour() *Receiver {
	if len(this.detours) == 0 {
		return nil
	}
	this.detourAccess.RLock()
	idx := 0
	detourLen := len(this.detours)
	if detourLen > 1 {
		idx = rand.Intn(detourLen)
	}
	rec := this.detours[idx]
	this.detourAccess.RUnlock()

	if rec.Expired() {
		this.detourAccess.Lock()
		detourLen := len(this.detours)
		this.detours[idx] = this.detours[detourLen-1]
		this.detours = this.detours[:detourLen-1]
		this.detourAccess.Unlock()
		return nil
	}

	return rec.Receiver
}

func (this *ReceiverManager) pickStdReceiver() *Receiver {
	receiverLen := len(this.receivers)

	receiverIdx := 0
	if receiverLen > 1 {
		receiverIdx = rand.Intn(receiverLen)
	}

	return this.receivers[receiverIdx]
}

func (this *ReceiverManager) PickReceiver() (v2net.Destination, *vmess.User) {
	rec := this.pickDetour()
	if rec == nil {
		rec = this.pickStdReceiver()
	}
	user := rec.PickUser()

	return rec.Destination, user
}
