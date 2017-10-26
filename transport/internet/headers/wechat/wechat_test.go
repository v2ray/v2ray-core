package wechat_test

import (
	"testing"

	"v2ray.com/core/common/buf"
	. "v2ray.com/core/transport/internet/headers/wechat"
	. "v2ray.com/ext/assert"
)

func TestUTPWrite(t *testing.T) {
	assert := With(t)

	video := VideoChat{}

	payload := buf.NewLocal(2048)
	payload.AppendSupplier(video.Write)

	assert(payload.Len(), Equals, video.Size())
}
