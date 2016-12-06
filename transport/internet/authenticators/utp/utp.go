package utp

import (
	"math/rand"

	"v2ray.com/core/common/alloc"
	"v2ray.com/core/common/loader"
	"v2ray.com/core/common/serial"
	"v2ray.com/core/transport/internet"
)

type UTP struct {
	header       byte
	extension    byte
	connectionId uint16
}

func (v *UTP) Overhead() int {
	return 4
}

func (v *UTP) Open(payload *alloc.Buffer) bool {
	payload.SliceFrom(v.Overhead())
	return true
}

func (v *UTP) Seal(payload *alloc.Buffer) {
	payload.PrependFunc(2, serial.WriteUint16(v.connectionId))
	payload.PrependBytes(v.header, v.extension)
}

type UTPFactory struct{}

func (v UTPFactory) Create(rawSettings interface{}) internet.Authenticator {
	return &UTP{
		header:       1,
		extension:    0,
		connectionId: uint16(rand.Intn(65536)),
	}
}

func init() {
	internet.RegisterAuthenticator(loader.GetType(new(Config)), UTPFactory{})
}
