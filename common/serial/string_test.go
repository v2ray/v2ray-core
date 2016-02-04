package serial_test

import (
	"testing"

	. "github.com/v2ray/v2ray-core/common/serial"
	v2testing "github.com/v2ray/v2ray-core/testing"
	"github.com/v2ray/v2ray-core/testing/assert"
)

type TestString struct {
	value string
}

func (this *TestString) String() string {
	return this.value
}

func TestNewStringSerial(t *testing.T) {
	v2testing.Current(t)

	testString := &TestString{value: "abcd"}
	assert.String(NewStringLiteral(testString)).Equals("abcd")
}
