package transport

import (
	"errors"
)

var (
	CorruptedPacket = errors.New("Packet is corrupted.")
)
