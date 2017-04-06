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
	assert.Bool(GetSeverity(err) == SeverityInfo).IsTrue()

	err = New("TestError2").Base(io.EOF)
	assert.Bool(GetSeverity(err) == SeverityInfo).IsTrue()

	err = New("TestError3").Base(io.EOF).AtWarning()
	assert.Bool(GetSeverity(err) == SeverityWarning).IsTrue()

	err = New("TestError4").Base(io.EOF).AtWarning()
	err = New("TestError5").Base(err)
	assert.Bool(GetSeverity(err) == SeverityWarning).IsTrue()
	assert.String(err.Error()).Contains("EOF")
}
