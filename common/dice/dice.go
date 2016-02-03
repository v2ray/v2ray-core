// Package dice contains common functions to generate random number.
// It also initialize math/rand with the time in seconds at launch time.

package dice

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

func init() {
	rand.Seed(time.Now().Unix())
}
