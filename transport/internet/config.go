package internet

import (
	"errors"

	"v2ray.com/core/common/loader"
	"v2ray.com/core/common/log"
	v2net "v2ray.com/core/common/net"
)

type ConfigCreator func() interface{}

var (
	globalNetworkConfigCreatorCache = make(map[v2net.Network]ConfigCreator)

	globalNetworkSettings []*NetworkSettings

	ErrUnconfiguredNetwork = errors.New("Network config creator not set.")
)

func RegisterNetworkConfigCreator(network v2net.Network, creator ConfigCreator) error {
	// TODO: check duplicate
	globalNetworkConfigCreatorCache[network] = creator
	return nil
}

func CreateNetworkConfig(network v2net.Network) (interface{}, error) {
	creator, ok := globalNetworkConfigCreatorCache[network]
	if !ok {
		log.Warning("Internet: Network config creator not found: ", network)
		return nil, ErrUnconfiguredNetwork
	}
	return creator(), nil
}

func (v *NetworkSettings) GetTypedSettings() (interface{}, error) {
	return v.Settings.GetInstance()
}

func (v *StreamConfig) GetEffectiveNetworkSettings() (interface{}, error) {
	for _, settings := range v.NetworkSettings {
		if settings.Network == v.Network {
			return settings.GetTypedSettings()
		}
	}
	for _, settings := range globalNetworkSettings {
		if settings.Network == v.Network {
			return settings.GetTypedSettings()
		}
	}
	return CreateNetworkConfig(v.Network)
}

func (v *StreamConfig) GetEffectiveSecuritySettings() (interface{}, error) {
	for _, settings := range v.SecuritySettings {
		if settings.Type == v.SecurityType {
			return settings.GetInstance()
		}
	}
	return loader.GetInstance(v.SecurityType)
}

func (v *StreamConfig) HasSecuritySettings() bool {
	return len(v.SecurityType) > 0
}

func ApplyGlobalNetworkSettings(settings []*NetworkSettings) error {
	globalNetworkSettings = settings
	return nil
}

func (v *ProxyConfig) HasTag() bool {
	return v != nil && len(v.Tag) > 0
}
