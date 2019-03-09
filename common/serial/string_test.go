package serial_test

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"

	. "v2ray.com/core/common/serial"
)

func TestToString(t *testing.T) {
	s := "a"
	data := []struct {
		Value  interface{}
		String string
	}{
		{Value: s, String: s},
		{Value: &s, String: s},
		{Value: errors.New("t"), String: "t"},
		{Value: []byte{'b', 'c'}, String: "[98 99]"},
	}

	for _, c := range data {
		if r := cmp.Diff(ToString(c.Value), c.String); r != "" {
			t.Error(r)
		}
	}
}

func TestConcat(t *testing.T) {
	testCases := []struct {
		Input  []interface{}
		Output string
	}{
		{
			Input: []interface{}{
				"a", "b",
			},
			Output: "ab",
		},
	}

	for _, testCase := range testCases {
		actual := Concat(testCase.Input...)
		if actual != testCase.Output {
			t.Error("Unexpected output: ", actual, " but want: ", testCase.Output)
		}
	}
}

func BenchmarkConcat(b *testing.B) {
	input := []interface{}{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k"}

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = Concat(input...)
	}
}
