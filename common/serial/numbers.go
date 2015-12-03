package serial

import (
	"strconv"
)

type Uint16 interface {
	Value() uint16
}

type Uint16Literal uint16

func (this Uint16Literal) String() string {
	return strconv.Itoa(int(this))
}

func (this Uint16Literal) Value() uint16 {
	return uint16(this)
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
