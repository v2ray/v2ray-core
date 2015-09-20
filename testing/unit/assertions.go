package unit

import (
	"testing"
)

// Assertion is an assertion library inspired by Truth.
// See http://google.github.io/truth/
type Assertion struct {
	t *testing.T
}

func Assert(t *testing.T) *Assertion {
	assert := new(Assertion)
	assert.t = t
	return assert
}

func (a *Assertion) Int64(value int64) *Int64Subject {
	return NewInt64Subject(NewSubject(a), value)
}

func (a *Assertion) Int(value int) *IntSubject {
	return NewIntSubject(NewSubject(a), value)
}

func (a *Assertion) Uint16(value uint16) *Uint16Subject {
	return NewUint16Subject(NewSubject(a), value)
}

func (a *Assertion) Byte(value byte) *ByteSubject {
	return NewByteSubject(NewSubject(a), value)
}

func (a *Assertion) Bytes(value []byte) *BytesSubject {
	return NewBytesSubject(NewSubject(a), value)
}

func (a *Assertion) String(value string) *StringSubject {
	return NewStringSubject(NewSubject(a), value)
}

func (a *Assertion) Error(value error) *ErrorSubject {
	return NewErrorSubject(NewSubject(a), value)
}

func (a *Assertion) Bool(value bool) *BoolSubject {
	return NewBoolSubject(NewSubject(a), value)
}

func (a *Assertion) Pointer(value interface{}) *PointerSubject {
	return NewPointerSubject(NewSubject(a), value)
}
