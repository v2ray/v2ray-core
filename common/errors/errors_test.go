package errors_test

import (
	"io"
	"testing"

	. "v2ray.com/core/common/errors"
	"v2ray.com/core/testing/assert"
)

func TestActionRequired(t *testing.T) {
	assert := assert.On(t)

	err := New("TestError")
	assert.Bool(IsActionRequired(err)).IsFalse()

	err = Base(io.EOF).Message("TestError2")
	assert.Bool(IsActionRequired(err)).IsFalse()

	err = Base(io.EOF).RequireUserAction().Message("TestError3")
	assert.Bool(IsActionRequired(err)).IsTrue()

	err = Base(io.EOF).RequireUserAction().Message("TestError4")
	err = Base(err).Message("TestError5")
	assert.Bool(IsActionRequired(err)).IsTrue()
	assert.String(err.Error()).Contains("EOF")
}
