package net

import (
	"strconv"
)

type Port uint16

func NewPort(port int) Port {
	return Port(uint16(port))
}

func (this Port) Value() uint16 {
	return uint16(this)
}

func (this Port) Bytes() []byte {
	return []byte{byte(this >> 8), byte(this)}
}

func (this Port) String() string {
	return strconv.Itoa(int(this))
}
