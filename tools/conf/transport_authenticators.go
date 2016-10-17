package conf

import (
	"v2ray.com/core/common/loader"
	"v2ray.com/core/transport/internet/authenticators/noop"
	"v2ray.com/core/transport/internet/authenticators/srtp"
	"v2ray.com/core/transport/internet/authenticators/utp"
)

type NoOpAuthenticator struct{}

func (NoOpAuthenticator) Build() (*loader.TypedSettings, error) {
	return loader.NewTypedSettings(new(noop.Config)), nil
}

type SRTPAuthenticator struct{}

func (SRTPAuthenticator) Build() (*loader.TypedSettings, error) {
	return loader.NewTypedSettings(new(srtp.Config)), nil
}

type UTPAuthenticator struct{}

func (UTPAuthenticator) Build() (*loader.TypedSettings, error) {
	return loader.NewTypedSettings(new(utp.Config)), nil
}
