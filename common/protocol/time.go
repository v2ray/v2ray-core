package protocol

import (
	"math/rand"
	"time"

	"github.com/v2ray/v2ray-core/common/serial"
)

type Timestamp int64

func (this Timestamp) Bytes() []byte {
	return serial.Int64ToBytes(int64(this))
}

type TimestampGenerator func() Timestamp

func NowTime() Timestamp {
	return Timestamp(time.Now().Unix())
}

func NewTimestampGenerator(base Timestamp, delta int) TimestampGenerator {
	return func() Timestamp {
		rangeInDelta := rand.Intn(delta*2) - delta
		return base + Timestamp(rangeInDelta)
	}
}
