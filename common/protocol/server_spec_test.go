package protocol_test

import (
	"testing"

	v2net "github.com/v2ray/v2ray-core/common/net"
	. "github.com/v2ray/v2ray-core/common/protocol"
	"github.com/v2ray/v2ray-core/common/uuid"
	"github.com/v2ray/v2ray-core/testing/assert"
)

func TestReceiverUser(t *testing.T) {
	assert := assert.On(t)

	id := NewID(uuid.New())
	alters := NewAlterIDs(id, 100)
	account := &VMessAccount{
		ID:       id,
		AlterIDs: alters,
	}
	user := NewUser(account, UserLevel(0), "")
	rec := NewServerSpec(v2net.TCPDestination(v2net.DomainAddress("v2ray.com"), 80), AlwaysValid(), user)
	assert.Bool(rec.HasUser(user)).IsTrue()

	id2 := NewID(uuid.New())
	alters2 := NewAlterIDs(id2, 100)
	account2 := &VMessAccount{
		ID:       id2,
		AlterIDs: alters2,
	}
	user2 := NewUser(account2, UserLevel(0), "")
	assert.Bool(rec.HasUser(user2)).IsFalse()

	rec.AddUser(user2)
	assert.Bool(rec.HasUser(user2)).IsTrue()
}
