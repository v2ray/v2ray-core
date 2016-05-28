package outbound_test

import (
	"testing"

	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/common/protocol"
	"github.com/v2ray/v2ray-core/common/uuid"
	. "github.com/v2ray/v2ray-core/proxy/vmess/outbound"
	"github.com/v2ray/v2ray-core/testing/assert"
)

func TestReceiverUser(t *testing.T) {
	assert := assert.On(t)

	id := protocol.NewID(uuid.New())
	alters := protocol.NewAlterIDs(id, 100)
	account := &protocol.VMessAccount{
		ID:       id,
		AlterIDs: alters,
	}
	user := protocol.NewUser(account, protocol.UserLevel(0), "")
	rec := NewReceiver(v2net.TCPDestination(v2net.DomainAddress("v2ray.com"), 80), user)
	assert.Bool(rec.HasUser(user)).IsTrue()
	assert.Int(len(rec.Accounts)).Equals(1)

	id2 := protocol.NewID(uuid.New())
	alters2 := protocol.NewAlterIDs(id2, 100)
	account2 := &protocol.VMessAccount{
		ID:       id2,
		AlterIDs: alters2,
	}
	user2 := protocol.NewUser(account2, protocol.UserLevel(0), "")
	assert.Bool(rec.HasUser(user2)).IsFalse()

	rec.AddUser(user2)
	assert.Bool(rec.HasUser(user2)).IsTrue()
	assert.Int(len(rec.Accounts)).Equals(2)
}
