package protocol

import (
	"math/rand"

	"github.com/v2ray/v2ray-core/common/protocol"
)

type RandomTimestampGenerator interface {
	Next() protocol.Timestamp
}

type RealRandomTimestampGenerator struct {
	base  protocol.Timestamp
	delta int
}

func NewRandomTimestampGenerator(base protocol.Timestamp, delta int) RandomTimestampGenerator {
	return &RealRandomTimestampGenerator{
		base:  base,
		delta: delta,
	}
}

func (this *RealRandomTimestampGenerator) Next() protocol.Timestamp {
	rangeInDelta := rand.Intn(this.delta*2) - this.delta
	return this.base + protocol.Timestamp(rangeInDelta)
}
