package internet

import (
	"v2ray.com/core/common/errors"
	"v2ray.com/core/common/serial"
)

type ConfigCreator func() interface{}

var (
	globalTransportConfigCreatorCache = make(map[TransportProtocol]ConfigCreator)
	globalTransportSettings           []*TransportConfig
)

func RegisterProtocolConfigCreator(protocol TransportProtocol, creator ConfigCreator) error {
	// TODO: check duplicate
	globalTransportConfigCreatorCache[protocol] = creator
	return nil
}

func CreateTransportConfig(protocol TransportProtocol) (interface{}, error) {
	creator, ok := globalTransportConfigCreatorCache[protocol]
	if !ok {
		return nil, errors.New("Internet: Unknown transport protocol: ", protocol)
	}
	return creator(), nil
}

func (v *TransportConfig) GetTypedSettings() (interface{}, error) {
	return v.Settings.GetInstance()
}

func (v *StreamConfig) GetEffectiveProtocol() TransportProtocol {
	if v == nil {
		return TransportProtocol_TCP
	}
	return v.Protocol
}

func (v *StreamConfig) GetEffectiveTransportSettings() (interface{}, error) {
	protocol := v.GetEffectiveProtocol()

	if v != nil {
		for _, settings := range v.TransportSettings {
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

func (v *StreamConfig) GetEffectiveSecuritySettings() (interface{}, error) {
	for _, settings := range v.SecuritySettings {
		if settings.Type == v.SecurityType {
			return settings.GetInstance()
		}
	}
	return serial.GetInstance(v.SecurityType)
}

func (v *StreamConfig) HasSecuritySettings() bool {
	return len(v.SecurityType) > 0
}

func ApplyGlobalTransportSettings(settings []*TransportConfig) error {
	globalTransportSettings = settings
	return nil
}

func (v *ProxyConfig) HasTag() bool {
	return v != nil && len(v.Tag) > 0
}
