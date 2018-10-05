package errors_test

import (
	"io"
	"testing"

	"v2ray.com/core/common/compare"
	. "v2ray.com/core/common/errors"
	"v2ray.com/core/common/log"
	. "v2ray.com/ext/assert"
)

func TestError(t *testing.T) {
	assert := With(t)

	err := New("TestError")
	assert(GetSeverity(err), Equals, log.Severity_Info)

	err = New("TestError2").Base(io.EOF)
	assert(GetSeverity(err), Equals, log.Severity_Info)

	err = New("TestError3").Base(io.EOF).AtWarning()
	assert(GetSeverity(err), Equals, log.Severity_Warning)

	err = New("TestError4").Base(io.EOF).AtWarning()
	err = New("TestError5").Base(err)
	assert(GetSeverity(err), Equals, log.Severity_Warning)
	assert(err.Error(), HasSubstring, "EOF")
}

type e struct{}

func TestErrorMessage(t *testing.T) {
	data := []struct {
		err error
		msg string
	}{
		{
			err: New("a").Base(New("b")).WithPathObj(e{}),
			msg: "v2ray.com/core/common/errors_test: a > b",
		},
		{
			err: New("a").Base(New("b").WithPathObj(e{})),
			msg: "a > v2ray.com/core/common/errors_test: b",
		},
	}

	for _, d := range data {
		if err := compare.StringEqualWithDetail(d.msg, d.err.Error()); err != nil {
			t.Fatal(err)
		}
	}
}
