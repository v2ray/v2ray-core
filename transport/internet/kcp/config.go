package kcp

import (
	"crypto/cipher"

	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/transport/internet"
)

// GetValue returns the value of MTU settings.
func (v *MTU) GetValue() uint32 {
	if v == nil {
		return 1350
	}
	return v.Value
}

// GetValue returns the value of TTI settings.
func (v *TTI) GetValue() uint32 {
	if v == nil {
		return 50
	}
	return v.Value
}

// GetValue returns the value of UplinkCapacity settings.
func (v *UplinkCapacity) GetValue() uint32 {
	if v == nil {
		return 5
	}
	return v.Value
}

// GetValue returns the value of DownlinkCapacity settings.
func (v *DownlinkCapacity) GetValue() uint32 {
	if v == nil {
		return 20
	}
	return v.Value
}

// GetSize returns the size of WriterBuffer in bytes.
func (v *WriteBuffer) GetSize() uint32 {
	if v == nil {
		return 2 * 1024 * 1024
	}
	return v.Size
}

// GetSize returns the size of ReadBuffer in bytes.
func (v *ReadBuffer) GetSize() uint32 {
	if v == nil {
		return 2 * 1024 * 1024
	}
	return v.Size
}

// GetSecurity returns the security settings.
func (v *Config) GetSecurity() (cipher.AEAD, error) {
	return NewSimpleAuthenticator(), nil
}

func (v *Config) GetPackerHeader() (internet.PacketHeader, error) {
	if v.HeaderConfig != nil {
		rawConfig, err := v.HeaderConfig.GetInstance()
		if err != nil {
			return nil, err
		}

		return internet.CreatePacketHeader(v.HeaderConfig.Type, rawConfig)
	}
	return nil, nil
}

func (v *Config) GetSendingInFlightSize() uint32 {
	size := v.UplinkCapacity.GetValue() * 1024 * 1024 / v.Mtu.GetValue() / (1000 / v.Tti.GetValue())
	if size < 8 {
		size = 8
	}
	return size
}

func (v *Config) GetSendingBufferSize() uint32 {
	return v.WriteBuffer.GetSize() / v.Mtu.GetValue()
}

func (v *Config) GetReceivingInFlightSize() uint32 {
	size := v.DownlinkCapacity.GetValue() * 1024 * 1024 / v.Mtu.GetValue() / (1000 / v.Tti.GetValue())
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
