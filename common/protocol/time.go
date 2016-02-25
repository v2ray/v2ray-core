package protocol

import (
	"math/rand"

	"github.com/v2ray/v2ray-core/common/serial"
)

type Timestamp int64

func (this Timestamp) Bytes() []byte {
	return serial.Int64Literal(this).Bytes()
}

type TimestampGenerator func() Timestamp

func NewTimestampGenerator(base Timestamp, delta int) TimestampGenerator {
	return func() Timestamp {
		rangeInDelta := rand.Intn(delta*2) - delta
		return base + Timestamp(rangeInDelta)
	}
}
