package mux_test

import (
	"testing"

	. "v2ray.com/core/app/proxyman/mux"
	"v2ray.com/core/testing/assert"
)

func TestSessionManagerAdd(t *testing.T) {
	assert := assert.On(t)

	m := NewSessionManager()

	s := &Session{}
	m.Allocate(s)
	assert.Uint16(s.ID).Equals(1)

	s = &Session{}
	m.Allocate(s)
	assert.Uint16(s.ID).Equals(2)

	s = &Session{
		ID: 4,
	}
	m.Add(s)
	assert.Uint16(s.ID).Equals(4)
}
