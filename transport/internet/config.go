package internet

import (
	"errors"

	"v2ray.com/core/common/log"
	v2net "v2ray.com/core/common/net"
	v2tls "v2ray.com/core/transport/internet/tls"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
)

type NetworkConfigCreator func() proto.Message

var (
	globalNetworkConfigCreatorCache  = make(map[v2net.Network]NetworkConfigCreator)
	globalSecurityConfigCreatorCache = make(map[SecurityType]NetworkConfigCreator)

	globalNetworkSettings []*NetworkSettings

	ErrUnconfiguredNetwork = errors.New("Network config creator not set.")
)

func RegisterNetworkConfigCreator(network v2net.Network, creator NetworkConfigCreator) error {
	// TODO: check duplicate
	globalNetworkConfigCreatorCache[network] = creator
	return nil
}

func CreateNetworkConfig(network v2net.Network) (proto.Message, error) {
	creator, ok := globalNetworkConfigCreatorCache[network]
	if !ok {
		log.Warning("Internet: Network config creator not found: ", network)
		return nil, ErrUnconfiguredNetwork
	}
	return creator(), nil
}

func RegisterSecurityConfigCreator(securityType SecurityType, creator NetworkConfigCreator) error {
	globalSecurityConfigCreatorCache[securityType] = creator
	return nil
}

func CreateSecurityConfig(securityType SecurityType) (proto.Message, error) {
	creator, ok := globalSecurityConfigCreatorCache[securityType]
	if !ok {
		log.Warning("Internet: Security config creator not found: ", securityType)
		return nil, ErrUnconfiguredNetwork
	}
	return creator(), nil
}

func (this *NetworkSettings) GetTypedSettings() (interface{}, error) {
	message, err := CreateNetworkConfig(this.Network)
	if err != nil {
		return nil, err
	}
	if err := ptypes.UnmarshalAny(this.Settings, message); err != nil {
		return nil, err
	}
	return message, nil
}

func (this *SecuritySettings) GetTypeSettings() (interface{}, error) {
	message, err := CreateSecurityConfig(this.Type)
	if err != nil {
		return nil, err
	}
	if err := ptypes.UnmarshalAny(this.Settings, message); err != nil {
		return nil, err
	}
	return message, nil
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
			return settings.GetTypeSettings()
		}
	}
	return CreateSecurityConfig(this.SecurityType)
}

func ApplyGlobalNetworkSettings(settings []*NetworkSettings) error {
	globalNetworkSettings = settings
	return nil
}

func init() {
	RegisterSecurityConfigCreator(SecurityType_TLS, func() proto.Message {
		return new(v2tls.Config)
	})
}
