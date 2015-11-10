package protocol

import (
	"testing"

	"github.com/v2ray/v2ray-core/testing/fuzzing"
)

const (
	Iterations = int(500000)
)

func TestReadAuthentication(t *testing.T) {
	for i := 0; i < Iterations; i++ {
		ReadAuthentication(fuzzing.RandomReader())
	}
}

func TestReadUserPassRequest(t *testing.T) {
	for i := 0; i < Iterations; i++ {
		ReadUserPassRequest(fuzzing.RandomReader())
	}
}

func TestReadRequest(t *testing.T) {
	for i := 0; i < Iterations; i++ {
		ReadRequest(fuzzing.RandomReader())
	}
}

func TestReadUDPRequest(t *testing.T) {
	for i := 0; i < Iterations; i++ {
		ReadUDPRequest(fuzzing.RandomBytes())
	}
}
