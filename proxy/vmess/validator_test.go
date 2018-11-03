package vmess_test

import (
	"testing"
	"time"

	"v2ray.com/core/common"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/common/serial"
	"v2ray.com/core/common/uuid"
	. "v2ray.com/core/proxy/vmess"
	. "v2ray.com/ext/assert"
)

func TestUserValidator(t *testing.T) {
	assert := With(t)

	hasher := protocol.DefaultIDHash
	v := NewTimedUserValidator(hasher)
	defer common.Close(v)

	toAccount := func(a *Account) protocol.Account {
		account, err := a.AsAccount()
		common.Must(err)
		return account
	}

	id := uuid.New()
	user := &protocol.MemoryUser{
		Email: "test",
		Account: toAccount(&Account{
			Id:      id.String(),
			AlterId: 8,
		}),
	}
	common.Must(v.Add(user))

	{
		testSmallLag := func(lag time.Duration) {
			ts := protocol.Timestamp(time.Now().Add(time.Second * lag).Unix())
			idHash := hasher(id.Bytes())
			common.Must2(serial.WriteUint64(idHash, uint64(ts)))
			userHash := idHash.Sum(nil)

			euser, ets, found := v.Get(userHash)
			assert(found, IsTrue)
			assert(euser.Email, Equals, user.Email)
			assert(int64(ets), Equals, int64(ts))
		}

		testSmallLag(0)
		testSmallLag(40)
		testSmallLag(-40)
		testSmallLag(80)
		testSmallLag(-80)
		testSmallLag(120)
		testSmallLag(-120)
	}

	{
		testBigLag := func(lag time.Duration) {
			ts := protocol.Timestamp(time.Now().Add(time.Second * lag).Unix())
			idHash := hasher(id.Bytes())
			common.Must2(serial.WriteUint64(idHash, uint64(ts)))
			userHash := idHash.Sum(nil)

			euser, _, found := v.Get(userHash)
			assert(found, IsFalse)
			assert(euser, IsNil)
		}

		testBigLag(121)
		testBigLag(-121)
		testBigLag(310)
		testBigLag(-310)
		testBigLag(500)
		testBigLag(-500)
	}

	assert(v.Remove(user.Email), IsTrue)
	assert(v.Remove(user.Email), IsFalse)
}
