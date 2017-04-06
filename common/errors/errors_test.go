package errors_test

import (
	"io"
	"testing"

	. "v2ray.com/core/common/errors"
	"v2ray.com/core/testing/assert"
)

func TestError(t *testing.T) {
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

func TestErrorMessage(t *testing.T) {
	assert := assert.On(t)

	data := []struct {
		err error
		msg string
	}{
		{
			err: New("a").Base(New("b")).Path("c", "d", "e"),
			msg: "c|d|e: a > b",
		},
		{
			err: New("a").Base(New("b").Path("c")).Path("d", "e"),
			msg: "d|e: a > c: b",
		},
	}

	for _, d := range data {
		assert.String(d.err.Error()).Equals(d.msg)
	}
}
