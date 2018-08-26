package vmess_test

import (
	"testing"
	"time"

	"v2ray.com/core/common"
	"v2ray.com/core/common/uuid"

	"v2ray.com/core/common/protocol"
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
		ts := protocol.Timestamp(time.Now().Unix())
		idHash := hasher(id.Bytes())
		idHash.Write(ts.Bytes(nil))
		userHash := idHash.Sum(nil)

		euser, ets, found := v.Get(userHash)
		assert(found, IsTrue)
		assert(euser.Email, Equals, user.Email)
		assert(int64(ets), Equals, int64(ts))
	}

	{
		ts := protocol.Timestamp(time.Now().Add(time.Second * 500).Unix())
		idHash := hasher(id.Bytes())
		idHash.Write(ts.Bytes(nil))
		userHash := idHash.Sum(nil)

		euser, _, found := v.Get(userHash)
		assert(found, IsFalse)
		assert(euser, IsNil)
	}

	assert(v.Remove(user.Email), IsTrue)
	assert(v.Remove(user.Email), IsFalse)
}
