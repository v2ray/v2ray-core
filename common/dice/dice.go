// Package dice contains common functions to generate random number.
// It also initialize math/rand with the time in seconds at launch time.
package dice // import "v2ray.com/core/common/dice"

import (
	"math/rand"
	"time"
)

// Roll returns a non-negative number between 0 (inclusive) and n (exclusive).
func Roll(n int) int {
	if n == 1 {
		return 0
	}
	return rand.Intn(n)
}

// Roll returns a non-negative number between 0 (inclusive) and n (exclusive).
func RollDeterministic(n int, seed int64) int {
	if n == 1 {
		return 0
	}
	return rand.New(rand.NewSource(seed)).Intn(n)
}

// RollUint16 returns a random uint16 value.
func RollUint16() uint16 {
	return uint16(rand.Intn(65536))
}

func RollUint64() uint64 {
	return rand.Uint64()
}

func NewDeterministicDice(seed int64) *deterministicDice {
	return &deterministicDice{rand.New(rand.NewSource(seed))}
}

type deterministicDice struct {
	*rand.Rand
}

func (dd *deterministicDice) Roll(n int) int {
	if n == 1 {
		return 0
	}
	return dd.Intn(n)
}

func init() {
	rand.Seed(time.Now().Unix())
}
