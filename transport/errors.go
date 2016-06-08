package transport

import (
	"errors"
)

var (
	ErrorCorruptedPacket = errors.New("Packet is corrupted.")
)
