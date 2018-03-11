package wechat_test

import (
	"context"
	"testing"

	"v2ray.com/core/common/buf"
	. "v2ray.com/core/transport/internet/headers/wechat"
	. "v2ray.com/ext/assert"
)

func TestUTPWrite(t *testing.T) {
	assert := With(t)

	videoRaw, err := NewVideoChat(context.Background(), &VideoConfig{})
	assert(err, IsNil)

	video := videoRaw.(*VideoChat)

	payload := buf.NewSize(2048)
	payload.AppendSupplier(video.Write)

	assert(payload.Len(), Equals, video.Size())
}
