package common_test

import (
	"errors"
	"testing"

	. "v2ray.com/core/common"
)

func TestMust(t *testing.T) {
	hasPanic := func(f func()) (ret bool) {
		defer func() {
			if r := recover(); r != nil {
				ret = true
			}
		}()
		f()
		return false
	}

	testCases := []struct {
		Input func()
		Panic bool
	}{
		{
			Panic: true,
			Input: func() { Must(func() error { return errors.New("test error") }()) },
		},
		{
			Panic: true,
			Input: func() { Must2(func() (int, error) { return 0, errors.New("test error") }()) },
		},
		{
			Panic: false,
			Input: func() { Must(func() error { return nil }()) },
		},
	}

	for idx, test := range testCases {
		if hasPanic(test.Input) != test.Panic {
			t.Error("test case #", idx, " expect panic ", test.Panic, " but actually not")
		}
	}
}
