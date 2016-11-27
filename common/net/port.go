package net

import (
	"errors"
	"strconv"

	"v2ray.com/core/common/serial"
)

var (
	// ErrInvalidPortRage indicates an error during port range parsing.
	ErrInvalidPortRange = errors.New("Invalid port range.")
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
func PortFromInt(val uint32) (Port, error) {
	if val > 65535 {
		return Port(0), ErrInvalidPortRange
	}
	return Port(val), nil
}

// PortFromString converts a string to a Port.
// @error when the string is not an integer or the integral value is a not a valid Port.
func PortFromString(s string) (Port, error) {
	val, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return Port(0), ErrInvalidPortRange
	}
	return PortFromInt(uint32(val))
}

// Value return the correspoding uint16 value of v Port.
func (v Port) Value() uint16 {
	return uint16(v)
}

// Bytes returns the correspoding bytes of v Port, in big endian order.
func (v Port) Bytes(b []byte) []byte {
	return serial.Uint16ToBytes(v.Value(), b)
}

// String returns the string presentation of v Port.
func (v Port) String() string {
	return serial.Uint16ToString(v.Value())
}

func (v PortRange) FromPort() Port {
	return Port(v.From)
}

func (v PortRange) ToPort() Port {
	return Port(v.To)
}

// Contains returns true if the given port is within the range of v PortRange.
func (v PortRange) Contains(port Port) bool {
	return v.FromPort() <= port && port <= v.ToPort()
}
