package user

import (
	"math/rand"
)

type RandomTimestampGenerator interface {
	Next() Timestamp
}

type RealRandomTimestampGenerator struct {
	base  Timestamp
	delta int
}

func NewRandomTimestampGenerator(base Timestamp, delta int) RandomTimestampGenerator {
	return &RealRandomTimestampGenerator{
		base:  base,
		delta: delta,
	}
}

func (this *RealRandomTimestampGenerator) Next() Timestamp {
	rangeInDelta := rand.Intn(this.delta*2) - this.delta
	return this.base + Timestamp(rangeInDelta)
}
