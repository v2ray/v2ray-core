package vmess_test

import (
	"testing"
	"time"

	"v2ray.com/core/common"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/common/serial"
	"v2ray.com/core/common/uuid"
	. "v2ray.com/core/proxy/vmess"
)

func TestUserValidator(t *testing.T) {
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
			if !found {
				t.Fatal("user not found")
			}
			if euser.Email != user.Email {
				t.Error("unexpected user email: ", euser.Email, " want ", user.Email)
			}
			if ets != ts {
				t.Error("unexpected timestamp: ", ets, " want ", ts)
			}
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
			if found || euser != nil {
				t.Error("unexpected user")
			}
		}

		testBigLag(121)
		testBigLag(-121)
		testBigLag(310)
		testBigLag(-310)
		testBigLag(500)
		testBigLag(-500)
	}

	if v := v.Remove(user.Email); !v {
		t.Error("unable to remove user")
	}
	if v := v.Remove(user.Email); v {
		t.Error("remove user twice")
	}
}
