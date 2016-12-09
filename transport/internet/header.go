package internet

import "v2ray.com/core/common"

type PacketHeader interface {
	Size() int
	Write([]byte) (int, error)
}

type PacketHeaderFactory interface {
	Create(interface{}) PacketHeader
}

var (
	headerCache = make(map[string]PacketHeaderFactory)
)

func RegisterPacketHeader(name string, factory PacketHeaderFactory) error {
	if _, found := headerCache[name]; found {
		return common.ErrDuplicatedName
	}
	headerCache[name] = factory
	return nil
}

func CreatePacketHeader(name string, config interface{}) (PacketHeader, error) {
	factory, found := headerCache[name]
	if !found {
		return nil, common.ErrObjectNotFound
	}
	return factory.Create(config), nil
}
