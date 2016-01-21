package dice

import (
	"math/rand"
)

func Roll(n int) int {
	if n == 1 {
		return 0
	}
	return rand.Intn(n)
}
