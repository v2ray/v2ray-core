package srtp

import (
	"math/rand"

	"github.com/v2ray/v2ray-core/common/alloc"
	"github.com/v2ray/v2ray-core/transport/internet"
)

type Config struct {
	Version     byte
	Padding     bool
	Extension   bool
	CSRCCount   byte
	Marker      bool
	PayloadType byte
}

type ObfuscatorSRTP struct {
	header uint16
	number uint16
}

func (this *ObfuscatorSRTP) Overhead() int {
	return 4
}

func (this *ObfuscatorSRTP) Open(payload *alloc.Buffer) bool {
	payload.SliceFrom(this.Overhead())
	return true
}

func (this *ObfuscatorSRTP) Seal(payload *alloc.Buffer) {
	this.number++
	payload.PrependUint16(this.number)
	payload.PrependUint16(this.header)
}

type ObfuscatorSRTPFactory struct {
}

func (this ObfuscatorSRTPFactory) Create(rawSettings internet.AuthenticatorConfig) internet.Authenticator {
	return &ObfuscatorSRTP{
		header: 0xB5E8,
		number: uint16(rand.Intn(65536)),
	}
}

func init() {
	internet.RegisterAuthenticator("srtp", ObfuscatorSRTPFactory{}, func() interface{} { return new(Config) })
}
