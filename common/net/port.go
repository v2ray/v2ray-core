package net

import (
	"errors"
	"strconv"

	"github.com/v2ray/v2ray-core/common/serial"
)

var (
	// ErrorInvalidPortRage indicates an error during port range parsing.
	ErrorInvalidPortRange = errors.New("Invalid port range.")
)

// Port represents a network port in TCP and UDP protocol.
type Port uint16

// PortFromBytes converts a byte array to a Port, assuming bytes are in big endian order.
// @unsafe Caller must ensure that the byte array has at least 2 elements.
func PortFromBytes(port []byte) Port {
	return Port(serial.BytesToUint16(port))
}

// PortFromInt converts an integer to a Port.
// @error when the integer is not positive or larger then 65535
func PortFromInt(v int) (Port, error) {
	if v <= 0 || v > 65535 {
		return Port(0), ErrorInvalidPortRange
	}
	return Port(v), nil
}

// PortFromString converts a string to a Port.
// @error when the string is not an integer or the integral value is a not a valid Port.
func PortFromString(s string) (Port, error) {
	v, err := strconv.Atoi(s)
	if err != nil {
		return Port(0), ErrorInvalidPortRange
	}
	return PortFromInt(v)
}

// Value return the correspoding uint16 value of this Port.
func (this Port) Value() uint16 {
	return uint16(this)
}

// Bytes returns the correspoding bytes of this Port, in big endian order.
func (this Port) Bytes(b []byte) []byte {
	return serial.Uint16ToBytes(this.Value(), b)
}

// String returns the string presentation of this Port.
func (this Port) String() string {
	return serial.Uint16ToString(this.Value())
}

// PortRange represents a range of ports.
type PortRange struct {
	From Port
	To   Port
}

// Contains returns true if the given port is within the range of this PortRange.
func (this PortRange) Contains(port Port) bool {
	return this.From <= port && port <= this.To
}
