package kcp

import (
	"v2ray.com/core/transport/internet"
)

type Config struct {
	Mtu              uint32 // Maximum transmission unit
	Tti              uint32
	UplinkCapacity   uint32
	DownlinkCapacity uint32
	Congestion       bool
	WriteBuffer      uint32
	ReadBuffer       uint32
	HeaderType       string
	HeaderConfig     internet.AuthenticatorConfig
}

func (this *Config) Apply() {
	effectiveConfig = *this
}

func (this *Config) GetAuthenticator() (internet.Authenticator, error) {
	auth := NewSimpleAuthenticator()
	if this.HeaderConfig != nil {
		header, err := internet.CreateAuthenticator(this.HeaderType, this.HeaderConfig)
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
	size := this.WriteBuffer / this.Mtu
	if size < this.GetSendingInFlightSize() {
		size = this.GetSendingInFlightSize()
	}
	return size
}

func (this *Config) GetReceivingWindowSize() uint32 {
	size := this.DownlinkCapacity * 1024 * 1024 / this.Mtu / (1000 / this.Tti) / 2
	if size < 8 {
		size = 8
	}
	return size
}

func (this *Config) GetReceivingBufferSize() uint32 {
	bufferSize := this.ReadBuffer / this.Mtu
	windowSize := this.DownlinkCapacity * 1024 * 1024 / this.Mtu / (1000 / this.Tti) / 2
	if windowSize < 8 {
		windowSize = 8
	}
	if bufferSize < windowSize {
		bufferSize = windowSize
	}
	return bufferSize
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
