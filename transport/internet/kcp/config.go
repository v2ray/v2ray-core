package kcp

import (
	"v2ray.com/core/transport/internet"
)

func (this *Config) Apply() {
	effectiveConfig = *this
}

func (this *Config) GetAuthenticator() (internet.Authenticator, error) {
	auth := NewSimpleAuthenticator()
	if this.HeaderConfig != nil {
		header, err := this.HeaderConfig.CreateAuthenticator()
		if err != nil {
			return nil, err
		}
		auth = internet.NewAuthenticatorChain(header, auth)
	}
	return auth, nil
}

func (this *Config) GetSendingInFlightSize() uint32 {
	size := this.UplinkCapacity * 1024 * 1024 / this.Mtu / (1000 / this.Tti) / 2
	if size < 8 {
		size = 8
	}
	return size
}

func (this *Config) GetSendingBufferSize() uint32 {
	return this.GetSendingInFlightSize() + this.WriteBuffer/this.Mtu
}

func (this *Config) GetReceivingInFlightSize() uint32 {
	size := this.DownlinkCapacity * 1024 * 1024 / this.Mtu / (1000 / this.Tti) / 2
	if size < 8 {
		size = 8
	}
	return size
}

func (this *Config) GetReceivingBufferSize() uint32 {
	return this.GetReceivingInFlightSize() + this.ReadBuffer/this.Mtu
}

func DefaultConfig() Config {
	return Config{
		Mtu:              1350,
		Tti:              50,
		UplinkCapacity:   5,
		DownlinkCapacity: 20,
		Congestion:       false,
		WriteBuffer:      1 * 1024 * 1024,
		ReadBuffer:       1 * 1024 * 1024,
	}
}

var (
	effectiveConfig = DefaultConfig()
)
