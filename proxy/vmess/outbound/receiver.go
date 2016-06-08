package outbound

import (
	"sync"
	"time"

	"github.com/v2ray/v2ray-core/common/dice"
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/common/protocol"
)

type Receiver struct {
	sync.RWMutex
	Destination v2net.Destination
	Accounts    []*protocol.User
}

func NewReceiver(dest v2net.Destination, users ...*protocol.User) *Receiver {
	return &Receiver{
		Destination: dest,
		Accounts:    users,
	}
}

func (this *Receiver) HasUser(user *protocol.User) bool {
	this.RLock()
	defer this.RUnlock()
	account := user.Account.(*protocol.VMessAccount)
	for _, u := range this.Accounts {
		// TODO: handle AlterIds difference.
		uAccount := u.Account.(*protocol.VMessAccount)
		if uAccount.ID.Equals(account.ID) {
			return true
		}
	}
	return false
}

func (this *Receiver) AddUser(user *protocol.User) {
	if this.HasUser(user) {
		return
	}
	this.Lock()
	this.Accounts = append(this.Accounts, user)
	this.Unlock()
}

func (this *Receiver) PickUser() *protocol.User {
	return this.Accounts[dice.Roll(len(this.Accounts))]
}

type ExpiringReceiver struct {
	*Receiver
	until time.Time
}

func (this *ExpiringReceiver) Expired() bool {
	return this.until.Before(time.Now())
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
			break
		}
	}

	this.detourAccess.RUnlock()
	if !destExists {
		expRec := &ExpiringReceiver{
			Receiver: rec,
			until:    time.Now().Add(time.Duration(availableMin-1) * time.Minute),
		}
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
	idx := dice.Roll(len(this.detours))
	rec := this.detours[idx]
	this.detourAccess.RUnlock()

	if rec.Expired() {
		this.detourAccess.Lock()
		detourLen := len(this.detours)
		if detourLen > idx && this.detours[idx].Expired() {
			this.detours[idx] = this.detours[detourLen-1]
			this.detours = this.detours[:detourLen-1]
		}
		this.detourAccess.Unlock()
		return nil
	}

	return rec.Receiver
}

func (this *ReceiverManager) pickStdReceiver() *Receiver {
	return this.receivers[dice.Roll(len(this.receivers))]
}

func (this *ReceiverManager) PickReceiver() *Receiver {
	rec := this.pickDetour()
	if rec == nil {
		rec = this.pickStdReceiver()
	}
	return rec
}
