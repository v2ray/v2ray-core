package mux_test

import (
	"testing"

	. "v2ray.com/core/app/proxyman/mux"
	. "v2ray.com/ext/assert"
)

func TestSessionManagerAdd(t *testing.T) {
	assert := With(t)

	m := NewSessionManager()

	s := m.Allocate()
	assert(s.ID, Equals, uint16(1))
	assert(m.Size(), Equals, 1)

	s = m.Allocate()
	assert(s.ID, Equals, uint16(2))
	assert(m.Size(), Equals, 2)

	s = &Session{
		ID: 4,
	}
	m.Add(s)
	assert(s.ID, Equals, uint16(4))
}

func TestSessionManagerClose(t *testing.T) {
	assert := With(t)

	m := NewSessionManager()
	s := m.Allocate()

	assert(m.CloseIfNoSession(), IsFalse)
	m.Remove(s.ID)
	assert(m.CloseIfNoSession(), IsTrue)
}
