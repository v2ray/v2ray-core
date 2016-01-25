package dice

import (
	"math/rand"
	"time"
)

func Roll(n int) int {
	if n == 1 {
		return 0
	}
	return rand.Intn(n)
}

func init() {
	rand.Seed(time.Now().Unix())
}
