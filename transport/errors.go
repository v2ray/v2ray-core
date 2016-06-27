package transport

import (
	"errors"
)

var (
	ErrCorruptedPacket = errors.New("Packet is corrupted.")
)
