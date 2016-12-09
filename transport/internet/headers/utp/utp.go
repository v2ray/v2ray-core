package utp

import (
	"math/rand"

	"v2ray.com/core/common/loader"
	"v2ray.com/core/common/serial"
	"v2ray.com/core/transport/internet"
)

type UTP struct {
	header       byte
	extension    byte
	connectionId uint16
}

func (v *UTP) Size() int {
	return 4
}

func (v *UTP) Write(b []byte) (int, error) {
	b = serial.Uint16ToBytes(v.connectionId, b[:0])
	b = append(b, v.header, v.extension)
	return 4, nil
}

type UTPFactory struct{}

func (v UTPFactory) Create(rawSettings interface{}) internet.PacketHeader {
	return &UTP{
		header:       1,
		extension:    0,
		connectionId: uint16(rand.Intn(65536)),
	}
}

func init() {
	internet.RegisterPacketHeader(loader.GetType(new(Config)), UTPFactory{})
}
