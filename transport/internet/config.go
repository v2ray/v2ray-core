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

func (this *NetworkSettings) GetTypedSettings() (interface{}, error) {
	return this.Settings.GetInstance()
}

func (this *StreamConfig) GetEffectiveNetworkSettings() (interface{}, error) {
	for _, settings := range this.NetworkSettings {
		if settings.Network == this.Network {
			return settings.GetTypedSettings()
		}
	}
	for _, settings := range globalNetworkSettings {
		if settings.Network == this.Network {
			return settings.GetTypedSettings()
		}
	}
	return CreateNetworkConfig(this.Network)
}

func (this *StreamConfig) GetEffectiveSecuritySettings() (interface{}, error) {
	for _, settings := range this.SecuritySettings {
		if settings.Type == this.SecurityType {
			return settings.GetInstance()
		}
	}
	return loader.GetInstance(this.SecurityType)
}

func (this *StreamConfig) HasSecuritySettings() bool {
	return len(this.SecurityType) > 0
}

func ApplyGlobalNetworkSettings(settings []*NetworkSettings) error {
	globalNetworkSettings = settings
	return nil
}
