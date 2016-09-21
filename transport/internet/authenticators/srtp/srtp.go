package srtp

import (
	"math/rand"

	"v2ray.com/core/common/alloc"
	"v2ray.com/core/transport/internet"
)

type SRTP struct {
	header uint16
	number uint16
}

func (this *SRTP) Overhead() int {
	return 4
}

func (this *SRTP) Open(payload *alloc.Buffer) bool {
	payload.SliceFrom(this.Overhead())
	return true
}

func (this *SRTP) Seal(payload *alloc.Buffer) {
	this.number++
	payload.PrependUint16(this.number)
	payload.PrependUint16(this.header)
}

type SRTPFactory struct {
}

func (this SRTPFactory) Create(rawSettings interface{}) internet.Authenticator {
	return &SRTP{
		header: 0xB5E8,
		number: uint16(rand.Intn(65536)),
	}
}

func init() {
	internet.RegisterAuthenticator("srtp", SRTPFactory{})
	internet.RegisterAuthenticatorConfig("srtp", func() interface{} { return &Config{} })
}
