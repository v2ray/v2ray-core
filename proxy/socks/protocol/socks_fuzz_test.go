package protocol

import (
	"crypto/rand"
	"testing"
)

const (
	Iterations = int(500000)
)

func TestReadAuthentication(t *testing.T) {
	for i := 0; i < Iterations; i++ {
		ReadAuthentication(rand.Reader)
	}
}

func TestReadUserPassRequest(t *testing.T) {
	for i := 0; i < Iterations; i++ {
		ReadUserPassRequest(rand.Reader)
	}
}

func TestReadRequest(t *testing.T) {
	for i := 0; i < Iterations; i++ {
		ReadRequest(rand.Reader)
	}
}
