package srtp

import (
	"math/rand"

	"v2ray.com/core/common/alloc"
	"v2ray.com/core/common/loader"
	"v2ray.com/core/transport/internet"
)

type SRTP struct {
	header uint16
	number uint16
}

func (v *SRTP) Overhead() int {
	return 4
}

func (v *SRTP) Open(payload *alloc.Buffer) bool {
	payload.SliceFrom(v.Overhead())
	return true
}

func (v *SRTP) Seal(payload *alloc.Buffer) {
	v.number++
	payload.PrependUint16(v.number)
	payload.PrependUint16(v.header)
}

type SRTPFactory struct {
}

func (v SRTPFactory) Create(rawSettings interface{}) internet.Authenticator {
	return &SRTP{
		header: 0xB5E8,
		number: uint16(rand.Intn(65536)),
	}
}

func init() {
	internet.RegisterAuthenticator(loader.GetType(new(Config)), SRTPFactory{})
}
