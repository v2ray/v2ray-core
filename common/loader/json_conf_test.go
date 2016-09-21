// +build json

package loader_test

import (
	"testing"

	. "v2ray.com/core/common/loader"
	"v2ray.com/core/testing/assert"
)

type TestConfigA struct {
	V int
}

type TestConfigB struct {
	S string
}

func TestCreatorCache(t *testing.T) {
	assert := assert.On(t)

	cache := ConfigCreatorCache{}
	creator1 := func() interface{} { return &TestConfigA{} }
	creator2 := func() interface{} { return &TestConfigB{} }
	cache.RegisterCreator("1", creator1)

	loader := NewJSONConfigLoader(cache, "test", "")
	rawA, err := loader.LoadWithID([]byte(`{"V": 2}`), "1")
	assert.Error(err).IsNil()
	instA := rawA.(*TestConfigA)
	assert.Int(instA.V).Equals(2)

	cache.RegisterCreator("2", creator2)
	rawB, err := loader.LoadWithID([]byte(`{"S": "a"}`), "2")
	assert.Error(err).IsNil()
	instB := rawB.(*TestConfigB)
	assert.String(instB.S).Equals("a")
}
