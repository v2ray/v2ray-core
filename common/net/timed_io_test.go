package net_test

import (
	"testing"

	. "v2ray.com/core/common/net"
	"v2ray.com/core/testing/assert"
)

func TestTimeOutSettings(t *testing.T) {
	assert := assert.On(t)

	reader := NewTimeOutReader(8, nil)
	assert.Int(reader.GetTimeOut()).Equals(8)
	reader.SetTimeOut(8) // no op
	assert.Int(reader.GetTimeOut()).Equals(8)
	reader.SetTimeOut(9)
	assert.Int(reader.GetTimeOut()).Equals(9)
}
