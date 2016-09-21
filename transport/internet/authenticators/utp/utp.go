package utp

import (
	"math/rand"

	"v2ray.com/core/common/alloc"
	"v2ray.com/core/transport/internet"
)

type UTP struct {
	header       byte
	extension    byte
	connectionId uint16
}

func (this *UTP) Overhead() int {
	return 4
}

func (this *UTP) Open(payload *alloc.Buffer) bool {
	payload.SliceFrom(this.Overhead())
	return true
}

func (this *UTP) Seal(payload *alloc.Buffer) {
	payload.PrependUint16(this.connectionId)
	payload.PrependBytes(this.header, this.extension)
}

type UTPFactory struct{}

func (this UTPFactory) Create(rawSettings interface{}) internet.Authenticator {
	return &UTP{
		header:       1,
		extension:    0,
		connectionId: uint16(rand.Intn(65536)),
	}
}

func init() {
	internet.RegisterAuthenticator("utp", UTPFactory{})
	internet.RegisterAuthenticatorConfig("utp", func() interface{} { return &Config{} })
}
