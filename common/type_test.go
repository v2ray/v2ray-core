package common_test

import (
	"context"
	"testing"

	. "v2ray.com/core/common"
	. "v2ray.com/ext/assert"
)

type TConfig struct {
	value int
}

type YConfig struct {
	value string
}

func TestObjectCreation(t *testing.T) {
	assert := With(t)

	var f = func(ctx context.Context, t interface{}) (interface{}, error) {
		return func() int {
			return t.(*TConfig).value
		}, nil
	}

	Must(RegisterConfig((*TConfig)(nil), f))
	err := RegisterConfig((*TConfig)(nil), f)
	assert(err, IsNotNil)

	g, err := CreateObject(context.Background(), &TConfig{value: 2})
	assert(err, IsNil)
	assert(g.(func() int)(), Equals, 2)

	_, err = CreateObject(context.Background(), &YConfig{value: "T"})
	assert(err, IsNotNil)
}
