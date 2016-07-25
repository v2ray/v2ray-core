package protocol_test

import (
	"testing"
	"time"

	v2net "github.com/v2ray/v2ray-core/common/net"
	. "github.com/v2ray/v2ray-core/common/protocol"
	"github.com/v2ray/v2ray-core/testing/assert"
)

type TestAccount struct {
	id int
}

func (this *TestAccount) Equals(account Account) bool {
	return account.(*TestAccount).id == this.id
}

func TestServerSpecUser(t *testing.T) {
	assert := assert.On(t)

	account := &TestAccount{
		id: 0,
	}
	user := NewUser(UserLevel(0), "")
	user.Account = account
	rec := NewServerSpec(v2net.TCPDestination(v2net.DomainAddress("v2ray.com"), 80), AlwaysValid(), user)
	assert.Bool(rec.HasUser(user)).IsTrue()

	account2 := &TestAccount{
		id: 1,
	}
	user2 := NewUser(UserLevel(0), "")
	user2.Account = account2
	assert.Bool(rec.HasUser(user2)).IsFalse()

	rec.AddUser(user2)
	assert.Bool(rec.HasUser(user2)).IsTrue()
}

func TestAlwaysValidStrategy(t *testing.T) {
	assert := assert.On(t)

	strategy := AlwaysValid()
	assert.Bool(strategy.IsValid()).IsTrue()
	strategy.Invalidate()
	assert.Bool(strategy.IsValid()).IsTrue()
}

func TestTimeoutValidStrategy(t *testing.T) {
	assert := assert.On(t)

	strategy := BeforeTime(time.Now().Add(2 * time.Second))
	assert.Bool(strategy.IsValid()).IsTrue()
	time.Sleep(3 * time.Second)
	assert.Bool(strategy.IsValid()).IsFalse()

	strategy = BeforeTime(time.Now().Add(2 * time.Second))
	strategy.Invalidate()
	assert.Bool(strategy.IsValid()).IsFalse()
}
