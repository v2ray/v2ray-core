package wechat_test

import (
	"testing"

	"v2ray.com/core/common/buf"
	"v2ray.com/core/testing/assert"
	. "v2ray.com/core/transport/internet/headers/wechat"
)

func TestUTPWrite(t *testing.T) {
	assert := assert.On(t)

	video := VideoChat{}

	payload := buf.NewLocal(2048)
	payload.AppendSupplier(video.Write)

	assert.Int(payload.Len()).Equals(video.Size())
}
