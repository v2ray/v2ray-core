package unit

import (
	"testing"
)

type Assertion struct {
	t *testing.T
}

func Assert(t *testing.T) *Assertion {
	assert := new(Assertion)
	assert.t = t
	return assert
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
