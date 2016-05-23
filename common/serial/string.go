package serial

import (
	"strings"
)

// An interface for any objects that has string presentation.
type String interface {
	String() string
}

type StringT string

func NewStringT(str String) StringT {
	return StringT(str.String())
}

func (this StringT) Contains(str String) bool {
	return strings.Contains(this.String(), str.String())
}

func (this StringT) String() string {
	return string(this)
}

func (this StringT) ToLower() StringT {
	return StringT(strings.ToLower(string(this)))
}

func (this StringT) ToUpper() StringT {
	return StringT(strings.ToUpper(string(this)))
}

func (this StringT) TrimSpace() StringT {
	return StringT(strings.TrimSpace(string(this)))
}
