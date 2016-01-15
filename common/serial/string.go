package serial

import (
	"strings"
)

// An interface for any objects that has string presentation.
type String interface {
	String() string
}

type StringLiteral string

func NewStringLiteral(str String) StringLiteral {
	return StringLiteral(str.String())
}

func (this StringLiteral) String() string {
	return string(this)
}

func (this StringLiteral) ToLower() StringLiteral {
	return StringLiteral(strings.ToLower(string(this)))
}

func (this StringLiteral) ToUpper() StringLiteral {
	return StringLiteral(strings.ToUpper(string(this)))
}

func (this StringLiteral) TrimSpace() StringLiteral {
	return StringLiteral(strings.TrimSpace(string(this)))
}
