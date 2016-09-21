package kcp

import (
	"v2ray.com/core/transport/internet"
)

func (this *MTU) GetValue() uint32 {
	if this == nil {
		return 1350
	}
	return this.Value
}

func (this *TTI) GetValue() uint32 {
	if this == nil {
		return 50
	}
	return this.Value
}

func (this *UplinkCapacity) GetValue() uint32 {
	if this == nil {
		return 5
	}
	return this.Value
}

func (this *DownlinkCapacity) GetValue() uint32 {
	if this == nil {
		return 20
	}
	return this.Value
}

func (this *WriteBuffer) GetSize() uint32 {
	if this == nil {
		return 1 * 1024 * 1024
	}
	return this.Size
}

func (this *ReadBuffer) GetSize() uint32 {
	if this == nil {
		return 1 * 1024 * 1024
	}
	return this.Size
}

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
	size := this.UplinkCapacity.GetValue() * 1024 * 1024 / this.Mtu.GetValue() / (1000 / this.Tti.GetValue()) / 2
	if size < 8 {
		size = 8
	}
	return size
}

func (this *Config) GetSendingBufferSize() uint32 {
	return this.GetSendingInFlightSize() + this.WriteBuffer.GetSize()/this.Mtu.GetValue()
}

func (this *Config) GetReceivingInFlightSize() uint32 {
	size := this.DownlinkCapacity.GetValue() * 1024 * 1024 / this.Mtu.GetValue() / (1000 / this.Tti.GetValue()) / 2
	if size < 8 {
		size = 8
	}
	return size
}

func (this *Config) GetReceivingBufferSize() uint32 {
	return this.GetReceivingInFlightSize() + this.ReadBuffer.GetSize()/this.Mtu.GetValue()
}

var (
	effectiveConfig Config
)
