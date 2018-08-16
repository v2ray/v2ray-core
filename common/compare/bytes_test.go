package compare_test

import (
	"testing"

	. "v2ray.com/core/common/compare"
)

func TestBytesEqual(t *testing.T) {
	testCases := []struct {
		Input1 []byte
		Input2 []byte
		Result bool
	}{
		{
			Input1: []byte{},
			Input2: []byte{1},
			Result: false,
		},
		{
			Input1: nil,
			Input2: []byte{},
			Result: true,
		},
		{
			Input1: []byte{1},
			Input2: []byte{1},
			Result: true,
		},
		{
			Input1: []byte{1, 2},
			Input2: []byte{1, 3},
			Result: false,
		},
	}

	for _, testCase := range testCases {
		cmp := BytesEqual(testCase.Input1, testCase.Input2)
		if cmp != testCase.Result {
			t.Errorf("unexpected result %v from %v", cmp, testCase)
		}
	}
}
