package srtp

import (
	"math/rand"

	"v2ray.com/core/common/loader"
	"v2ray.com/core/common/serial"
	"v2ray.com/core/transport/internet"
)

type SRTP struct {
	header uint16
	number uint16
}

func (v *SRTP) Size() int {
	return 4
}

func (v *SRTP) Write(b []byte) (int, error) {
	v.number++
	b = serial.Uint16ToBytes(v.number, b[:0])
	b = serial.Uint16ToBytes(v.number, b)
	return 4, nil
}

type SRTPFactory struct {
}

func (v SRTPFactory) Create(rawSettings interface{}) internet.PacketHeader {
	return &SRTP{
		header: 0xB5E8,
		number: uint16(rand.Intn(65536)),
	}
}

func init() {
	internet.RegisterPacketHeader(loader.GetType(new(Config)), SRTPFactory{})
}
