package outbound_test

import (
	"testing"

	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/common/uuid"
	"github.com/v2ray/v2ray-core/proxy/vmess"
	. "github.com/v2ray/v2ray-core/proxy/vmess/outbound"
	v2testing "github.com/v2ray/v2ray-core/testing"
	"github.com/v2ray/v2ray-core/testing/assert"
)

func TestReceiverUser(t *testing.T) {
	v2testing.Current(t)

	id := vmess.NewID(uuid.New())
	user := vmess.NewUser(id, vmess.UserLevel(0), 100)
	rec := NewReceiver(v2net.TCPDestination(v2net.DomainAddress("v2ray.com"), 80), user)
	assert.Bool(rec.HasUser(user)).IsTrue()
	assert.Int(len(rec.Accounts)).Equals(1)

	id2 := vmess.NewID(uuid.New())
	user2 := vmess.NewUser(id2, vmess.UserLevel(0), 100)
	assert.Bool(rec.HasUser(user2)).IsFalse()

	rec.AddUser(user2)
	assert.Bool(rec.HasUser(user2)).IsTrue()
	assert.Int(len(rec.Accounts)).Equals(2)
}
