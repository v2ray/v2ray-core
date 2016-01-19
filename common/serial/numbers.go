package serial

import (
	"strconv"
)

type Uint16 interface {
	Value() uint16
}

type Uint16Literal uint16

func ParseUint16(data []byte) Uint16Literal {
	switch len(data) {
	case 0:
		return Uint16Literal(0)
	case 1:
		return Uint16Literal(uint16(data[0]))
	default:
		return Uint16Literal(uint16(data[0])<<8 + uint16(data[1]))
	}
}

func (this Uint16Literal) String() string {
	return strconv.Itoa(int(this))
}

func (this Uint16Literal) Value() uint16 {
	return uint16(this)
}

func (this Uint16Literal) Bytes() []byte {
	return []byte{byte(this >> 8), byte(this)}
}

type Int interface {
	Value() int
}

type IntLiteral int

func (this IntLiteral) String() string {
	return strconv.Itoa(int(this))
}

func (this IntLiteral) Value() int {
	return int(this)
}

type Int64Literal int64

func (this Int64Literal) String() string {
	return strconv.FormatInt(this.Value(), 10)
}

func (this Int64Literal) Value() int64 {
	return int64(this)
}

func (this Int64Literal) Bytes() []byte {
	value := this.Value()
	return []byte{
		byte(value >> 56),
		byte(value >> 48),
		byte(value >> 40),
		byte(value >> 32),
		byte(value >> 24),
		byte(value >> 16),
		byte(value >> 8),
		byte(value),
	}
}
