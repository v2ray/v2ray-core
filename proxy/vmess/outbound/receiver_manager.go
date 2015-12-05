package outbound

import (
	"math/rand"

	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/proxy/vmess/config"
)

type ReceiverManager struct {
	receivers []*config.Receiver
}

func NewReceiverManager(receivers []*config.Receiver) *ReceiverManager {
	return &ReceiverManager{
		receivers: receivers,
	}
}

func (this *ReceiverManager) PickReceiver() (v2net.Address, config.User) {
	receiverLen := len(this.receivers)
	receiverIdx := 0
	if receiverLen > 1 {
		receiverIdx = rand.Intn(receiverLen)
	}

	receiver := this.receivers[receiverIdx]

	userLen := len(receiver.Accounts)
	userIdx := 0
	if userLen > 1 {
		userIdx = rand.Intn(userLen)
	}
	user := receiver.Accounts[userIdx]
	return receiver.Address, user
}
