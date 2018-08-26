package protocol_test

import (
	"testing"
	"time"

	"v2ray.com/core/common/net"
	. "v2ray.com/core/common/protocol"
	"v2ray.com/core/common/uuid"
	"v2ray.com/core/proxy/vmess"
	. "v2ray.com/ext/assert"
)

func TestAlwaysValidStrategy(t *testing.T) {
	assert := With(t)

	strategy := AlwaysValid()
	assert(strategy.IsValid(), IsTrue)
	strategy.Invalidate()
	assert(strategy.IsValid(), IsTrue)
}

func TestTimeoutValidStrategy(t *testing.T) {
	assert := With(t)

	strategy := BeforeTime(time.Now().Add(2 * time.Second))
	assert(strategy.IsValid(), IsTrue)
	time.Sleep(3 * time.Second)
	assert(strategy.IsValid(), IsFalse)

	strategy = BeforeTime(time.Now().Add(2 * time.Second))
	strategy.Invalidate()
	assert(strategy.IsValid(), IsFalse)
}

func TestUserInServerSpec(t *testing.T) {
	assert := With(t)

	uuid1 := uuid.New()
	uuid2 := uuid.New()

	spec := NewServerSpec(net.Destination{}, AlwaysValid(), &MemoryUser{
		Email:   "test1@v2ray.com",
		Account: &vmess.Account{Id: uuid1.String()},
	})
	assert(spec.HasUser(&MemoryUser{
		Email:   "test1@v2ray.com",
		Account: &vmess.Account{Id: uuid2.String()},
	}), IsFalse)

	spec.AddUser(&MemoryUser{Email: "test2@v2ray.com"})
	assert(spec.HasUser(&MemoryUser{
		Email:   "test1@v2ray.com",
		Account: &vmess.Account{Id: uuid1.String()},
	}), IsTrue)
}

func TestPickUser(t *testing.T) {
	assert := With(t)

	spec := NewServerSpec(net.Destination{}, AlwaysValid(), &MemoryUser{Email: "test1@v2ray.com"}, &MemoryUser{Email: "test2@v2ray.com"}, &MemoryUser{Email: "test3@v2ray.com"})
	user := spec.PickUser()
	assert(user.Email, HasSuffix, "@v2ray.com")
}
