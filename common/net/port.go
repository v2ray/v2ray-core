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
func PortFromInt(v uint32) (Port, error) {
	if v > 65535 {
		return Port(0), ErrInvalidPortRange
	}
	return Port(v), nil
}

// PortFromString converts a string to a Port.
// @error when the string is not an integer or the integral value is a not a valid Port.
func PortFromString(s string) (Port, error) {
	v, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return Port(0), ErrInvalidPortRange
	}
	return PortFromInt(uint32(v))
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

func (this PortRange) FromPort() Port {
	return Port(this.From)
}

func (this PortRange) ToPort() Port {
	return Port(this.To)
}

// Contains returns true if the given port is within the range of this PortRange.
func (this PortRange) Contains(port Port) bool {
	return this.FromPort() <= port && port <= this.ToPort()
}
