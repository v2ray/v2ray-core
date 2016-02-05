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

type Port serial.Uint16Literal

func PortFromBytes(port []byte) Port {
	return Port(serial.BytesLiteral(port).Uint16Value())
}

func PortFromInt(v int) (Port, error) {
	if v <= 0 || v > 65535 {
		return Port(0), ErrorInvalidPortRange
	}
	return Port(v), nil
}

func PortFromString(s string) (Port, error) {
	v, err := strconv.Atoi(s)
	if err != nil {
		return Port(0), ErrorInvalidPortRange
	}
	return PortFromInt(v)
}

func (this Port) Value() uint16 {
	return uint16(this)
}

func (this Port) Bytes() []byte {
	return serial.Uint16Literal(this).Bytes()
}

func (this Port) String() string {
	return serial.Uint16Literal(this).String()
}

type PortRange struct {
	From Port
	To   Port
}

func (this PortRange) Contains(port Port) bool {
	return this.From <= port && port <= this.To
}
