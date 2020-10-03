package encoding_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/common/uuid"
	"v2ray.com/core/proxy/vless"
	. "v2ray.com/core/proxy/vless/encoding"
)

func toAccount(a *vless.Account) protocol.Account {
	account, err := a.AsAccount()
	common.Must(err)
	return account
}

func TestRequestSerialization(t *testing.T) {
	user := &protocol.MemoryUser{
		Level: 0,
		Email: "test@v2fly.org",
	}
	id := uuid.New()
	account := &vless.Account{
		Id: id.String(),
	}
	user.Account = toAccount(account)

	expectedRequest := &protocol.RequestHeader{
		Version: Version,
		User:    user,
		Command: protocol.RequestCommandTCP,
		Address: net.DomainAddress("www.v2fly.org"),
		Port:    net.Port(443),
	}
	expectedAddons := &Addons{}

	buffer := buf.StackNew()
	common.Must(EncodeRequestHeader(&buffer, expectedRequest, expectedAddons))

	Validator := new(vless.Validator)
	Validator.Add(user)

	actualRequest, actualAddons, err, _ := DecodeRequestHeader(false, nil, &buffer, Validator)
	common.Must(err)

	if r := cmp.Diff(actualRequest, expectedRequest, cmp.AllowUnexported(protocol.ID{})); r != "" {
		t.Error(r)
	}
	if r := cmp.Diff(actualAddons, expectedAddons); r != "" {
		t.Error(r)
	}
}

func TestInvalidRequest(t *testing.T) {
	user := &protocol.MemoryUser{
		Level: 0,
		Email: "test@v2fly.org",
	}
	id := uuid.New()
	account := &vless.Account{
		Id: id.String(),
	}
	user.Account = toAccount(account)

	expectedRequest := &protocol.RequestHeader{
		Version: Version,
		User:    user,
		Command: protocol.RequestCommand(100),
		Address: net.DomainAddress("www.v2fly.org"),
		Port:    net.Port(443),
	}
	expectedAddons := &Addons{}

	buffer := buf.StackNew()
	common.Must(EncodeRequestHeader(&buffer, expectedRequest, expectedAddons))

	Validator := new(vless.Validator)
	Validator.Add(user)

	_, _, err, _ := DecodeRequestHeader(false, nil, &buffer, Validator)
	if err == nil {
		t.Error("nil error")
	}
}

func TestMuxRequest(t *testing.T) {
	user := &protocol.MemoryUser{
		Level: 0,
		Email: "test@v2fly.org",
	}
	id := uuid.New()
	account := &vless.Account{
		Id: id.String(),
	}
	user.Account = toAccount(account)

	expectedRequest := &protocol.RequestHeader{
		Version: Version,
		User:    user,
		Command: protocol.RequestCommandMux,
		Address: net.DomainAddress("v1.mux.cool"),
	}
	expectedAddons := &Addons{}

	buffer := buf.StackNew()
	common.Must(EncodeRequestHeader(&buffer, expectedRequest, expectedAddons))

	Validator := new(vless.Validator)
	Validator.Add(user)

	actualRequest, actualAddons, err, _ := DecodeRequestHeader(false, nil, &buffer, Validator)
	common.Must(err)

	if r := cmp.Diff(actualRequest, expectedRequest, cmp.AllowUnexported(protocol.ID{})); r != "" {
		t.Error(r)
	}
	if r := cmp.Diff(actualAddons, expectedAddons); r != "" {
		t.Error(r)
	}
}
