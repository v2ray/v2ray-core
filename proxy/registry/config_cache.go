package registry

import (
	"v2ray.com/core/common/loader"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
)

var (
	inboundConfigCreatorCache  = loader.ConfigCreatorCache{}
	outboundConfigCreatorCache = loader.ConfigCreatorCache{}
)

func RegisterInboundConfig(protocol string, creator loader.ConfigCreator) error {
	return inboundConfigCreatorCache.RegisterCreator(protocol, creator)
}

func RegisterOutboundConfig(protocol string, creator loader.ConfigCreator) error {
	return outboundConfigCreatorCache.RegisterCreator(protocol, creator)
}

func MarshalInboundConfig(protocol string, settings *any.Any) (interface{}, error) {
	config, err := inboundConfigCreatorCache.CreateConfig(protocol)
	if err != nil {
		return nil, err
	}
	if err := ptypes.UnmarshalAny(settings, config.(proto.Message)); err != nil {
		return nil, err
	}
	return config, nil
}

func MarshalOutboundConfig(protocol string, settings *any.Any) (interface{}, error) {
	config, err := outboundConfigCreatorCache.CreateConfig(protocol)
	if err != nil {
		return nil, err
	}
	if err := ptypes.UnmarshalAny(settings, config.(proto.Message)); err != nil {
		return nil, err
	}
	return config, nil
}
