package kcp

import (
	"crypto/cipher"

	"v2ray.com/core/transport/internet"
)

const protocolName = "mkcp"

// GetMTUValue returns the value of MTU settings.
func (c *Config) GetMTUValue() uint32 {
	if c == nil || c.Mtu == nil {
		return 1350
	}
	return c.Mtu.Value
}

// GetTTIValue returns the value of TTI settings.
func (c *Config) GetTTIValue() uint32 {
	if c == nil || c.Tti == nil {
		return 50
	}
	return c.Tti.Value
}

// GetUplinkCapacityValue returns the value of UplinkCapacity settings.
func (c *Config) GetUplinkCapacityValue() uint32 {
	if c == nil || c.UplinkCapacity == nil {
		return 5
	}
	return c.UplinkCapacity.Value
}

// GetDownlinkCapacityValue returns the value of DownlinkCapacity settings.
func (c *Config) GetDownlinkCapacityValue() uint32 {
	if c == nil || c.DownlinkCapacity == nil {
		return 20
	}
	return c.DownlinkCapacity.Value
}

// GetWriteBufferSize returns the size of WriterBuffer in bytes.
func (c *Config) GetWriteBufferSize() uint32 {
	if c == nil || c.WriteBuffer == nil {
		return 2 * 1024 * 1024
	}
	return c.WriteBuffer.Size
}

// GetReadBufferSize returns the size of ReadBuffer in bytes.
func (c *Config) GetReadBufferSize() uint32 {
	if c == nil || c.ReadBuffer == nil {
		return 2 * 1024 * 1024
	}
	return c.ReadBuffer.Size
}

// GetSecurity returns the security settings.
func (*Config) GetSecurity() (cipher.AEAD, error) {
	return NewSimpleAuthenticator(), nil
}

func (c *Config) GetPackerHeader() (internet.PacketHeader, error) {
	if c.HeaderConfig != nil {
		rawConfig, err := c.HeaderConfig.GetInstance()
		if err != nil {
			return nil, err
		}

		return internet.CreatePacketHeader(rawConfig)
	}
	return nil, nil
}

func (c *Config) GetSendingInFlightSize() uint32 {
	size := c.GetUplinkCapacityValue() * 1024 * 1024 / c.GetMTUValue() / (1000 / c.GetTTIValue())
	if size < 8 {
		size = 8
	}
	return size
}

func (c *Config) GetSendingBufferSize() uint32 {
	return c.GetWriteBufferSize() / c.GetMTUValue()
}

func (c *Config) GetReceivingInFlightSize() uint32 {
	size := c.GetDownlinkCapacityValue() * 1024 * 1024 / c.GetMTUValue() / (1000 / c.GetTTIValue())
	if size < 8 {
		size = 8
	}
	return size
}

func (c *Config) GetReceivingBufferSize() uint32 {
	return c.GetReadBufferSize() / c.GetMTUValue()
}
