package kcp

import (
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/transport/internet"
)

func (v *MTU) GetValue() uint32 {
	if v == nil {
		return 1350
	}
	return v.Value
}

func (v *TTI) GetValue() uint32 {
	if v == nil {
		return 50
	}
	return v.Value
}

func (v *UplinkCapacity) GetValue() uint32 {
	if v == nil {
		return 5
	}
	return v.Value
}

func (v *DownlinkCapacity) GetValue() uint32 {
	if v == nil {
		return 20
	}
	return v.Value
}

func (v *WriteBuffer) GetSize() uint32 {
	if v == nil {
		return 1 * 1024 * 1024
	}
	return v.Size
}

func (v *ReadBuffer) GetSize() uint32 {
	if v == nil {
		return 1 * 1024 * 1024
	}
	return v.Size
}

func (v *Config) GetAuthenticator() (internet.Authenticator, error) {
	auth := NewSimpleAuthenticator()
	if v.HeaderConfig != nil {
		rawConfig, err := v.HeaderConfig.GetInstance()
		if err != nil {
			return nil, err
		}

		header, err := internet.CreateAuthenticator(v.HeaderConfig.Type, rawConfig)
		if err != nil {
			return nil, err
		}
		auth = internet.NewAuthenticatorChain(header, auth)
	}
	return auth, nil
}

func (v *Config) GetSendingInFlightSize() uint32 {
	size := v.UplinkCapacity.GetValue() * 1024 * 1024 / v.Mtu.GetValue() / (1000 / v.Tti.GetValue()) / 2
	if size < 8 {
		size = 8
	}
	return size
}

func (v *Config) GetSendingBufferSize() uint32 {
	return v.WriteBuffer.GetSize() / v.Mtu.GetValue()
}

func (v *Config) GetReceivingInFlightSize() uint32 {
	size := v.DownlinkCapacity.GetValue() * 1024 * 1024 / v.Mtu.GetValue() / (1000 / v.Tti.GetValue()) / 2
	if size < 8 {
		size = 8
	}
	return size
}

func (v *Config) GetReceivingBufferSize() uint32 {
	return v.ReadBuffer.GetSize() / v.Mtu.GetValue()
}

func (o *ConnectionReuse) IsEnabled() bool {
	if o == nil {
		return true
	}
	return o.Enable
}

func init() {
	internet.RegisterNetworkConfigCreator(v2net.Network_KCP, func() interface{} {
		return new(Config)
	})
}
