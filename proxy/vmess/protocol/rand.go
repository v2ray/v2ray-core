package protocol

import (
	"math/rand"
)

type RandomInt64InRange func(base int64, delta int) int64

func GenerateRandomInt64InRange(base int64, delta int) int64 {
	rangeInDelta := rand.Intn(delta*2) - delta
	return base + int64(rangeInDelta)
}
