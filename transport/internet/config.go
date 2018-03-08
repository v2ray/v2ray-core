package internet

import "v2ray.com/core/common/serial"

type ConfigCreator func() interface{}

var (
	globalTransportConfigCreatorCache = make(map[TransportProtocol]ConfigCreator)
	globalTransportSettings           []*TransportConfig
)

func RegisterProtocolConfigCreator(protocol TransportProtocol, creator ConfigCreator) error {
	if _, found := globalTransportConfigCreatorCache[protocol]; found {
		return newError("protocol: " + TransportProtocol_name[int32(protocol)]+ " is already registered").AtError()
	}
	globalTransportConfigCreatorCache[protocol] = creator
	return nil
}

func CreateTransportConfig(protocol TransportProtocol) (interface{}, error) {
	creator, ok := globalTransportConfigCreatorCache[protocol]
	if !ok {
		return nil, newError("unknown transport protocol: ", protocol)
	}
	return creator(), nil
}

func (c *TransportConfig) GetTypedSettings() (interface{}, error) {
	return c.Settings.GetInstance()
}

func (c *StreamConfig) GetEffectiveProtocol() TransportProtocol {
	if c == nil {
		return TransportProtocol_TCP
	}
	return c.Protocol
}

func (c *StreamConfig) GetEffectiveTransportSettings() (interface{}, error) {
	protocol := c.GetEffectiveProtocol()

	if c != nil {
		for _, settings := range c.TransportSettings {
			if settings.Protocol == protocol {
				return settings.GetTypedSettings()
			}
		}
	}

	for _, settings := range globalTransportSettings {
		if settings.Protocol == protocol {
			return settings.GetTypedSettings()
		}
	}
	return CreateTransportConfig(protocol)
}

func (c *StreamConfig) GetTransportSettingsFor(protocol TransportProtocol) (interface{}, error) {
	if c != nil {
		for _, settings := range c.TransportSettings {
			if settings.Protocol == protocol {
				return settings.GetTypedSettings()
			}
		}
	}

	for _, settings := range globalTransportSettings {
		if settings.Protocol == protocol {
			return settings.GetTypedSettings()
		}
	}

	return CreateTransportConfig(protocol)
}

func (c *StreamConfig) GetEffectiveSecuritySettings() (interface{}, error) {
	for _, settings := range c.SecuritySettings {
		if settings.Type == c.SecurityType {
			return settings.GetInstance()
		}
	}
	return serial.GetInstance(c.SecurityType)
}

func (c *StreamConfig) HasSecuritySettings() bool {
	return len(c.SecurityType) > 0
}

func ApplyGlobalTransportSettings(settings []*TransportConfig) error {
	globalTransportSettings = settings
	return nil
}

func (c *ProxyConfig) HasTag() bool {
	return c != nil && len(c.Tag) > 0
}
