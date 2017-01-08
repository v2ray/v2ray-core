package kcp

import (
	"crypto/cipher"

	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/transport/internet"
)

// GetMTUValue returns the value of MTU settings.
func (v *Config) GetMTUValue() uint32 {
	if v == nil || v.Mtu == nil {
		return 1350
	}
	return v.Mtu.Value
}

// GetTTIValue returns the value of TTI settings.
func (v *Config) GetTTIValue() uint32 {
	if v == nil || v.Tti == nil {
		return 50
	}
	return v.Tti.Value
}

// GetUplinkCapacityValue returns the value of UplinkCapacity settings.
func (v *Config) GetUplinkCapacityValue() uint32 {
	if v == nil || v.UplinkCapacity == nil {
		return 5
	}
	return v.UplinkCapacity.Value
}

// GetDownlinkCapacityValue returns the value of DownlinkCapacity settings.
func (v *Config) GetDownlinkCapacityValue() uint32 {
	if v == nil || v.DownlinkCapacity == nil {
		return 20
	}
	return v.DownlinkCapacity.Value
}

// GetWriteBufferSize returns the size of WriterBuffer in bytes.
func (v *Config) GetWriteBufferSize() uint32 {
	if v == nil || v.WriteBuffer == nil {
		return 2 * 1024 * 1024
	}
	return v.WriteBuffer.Size
}

// GetReadBufferSize returns the size of ReadBuffer in bytes.
func (v *Config) GetReadBufferSize() uint32 {
	if v == nil || v.ReadBuffer == nil {
		return 2 * 1024 * 1024
	}
	return v.ReadBuffer.Size
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
	size := v.GetUplinkCapacityValue() * 1024 * 1024 / v.GetMTUValue() / (1000 / v.GetTTIValue())
	if size < 8 {
		size = 8
	}
	return size
}

func (v *Config) GetSendingBufferSize() uint32 {
	return v.GetWriteBufferSize() / v.GetMTUValue()
}

func (v *Config) GetReceivingInFlightSize() uint32 {
	size := v.GetDownlinkCapacityValue() * 1024 * 1024 / v.GetMTUValue() / (1000 / v.GetTTIValue())
	if size < 8 {
		size = 8
	}
	return size
}

func (v *Config) GetReceivingBufferSize() uint32 {
	return v.GetReadBufferSize() / v.GetMTUValue()
}

func (v *Config) IsConnectionReuse() bool {
	if v == nil || v.ConnectionReuse == nil {
		return true
	}
	return v.ConnectionReuse.Enable
}

func init() {
	internet.RegisterNetworkConfigCreator(v2net.Network_KCP, func() interface{} {
		return new(Config)
	})
}
