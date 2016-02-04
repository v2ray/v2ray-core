package net_test

import (
	"testing"

	. "github.com/v2ray/v2ray-core/common/net"
	v2testing "github.com/v2ray/v2ray-core/testing"
	"github.com/v2ray/v2ray-core/testing/assert"
)

func TestTimeOutSettings(t *testing.T) {
	v2testing.Current(t)

	reader := NewTimeOutReader(8, nil)
	assert.Int(reader.GetTimeOut()).Equals(8)
	reader.SetTimeOut(8) // no op
	assert.Int(reader.GetTimeOut()).Equals(8)
	reader.SetTimeOut(9)
	assert.Int(reader.GetTimeOut()).Equals(9)
}
