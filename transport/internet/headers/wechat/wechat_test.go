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

	payload := buf.New()
	video.Serialize(payload.Extend(video.Size()))

	assert(payload.Len(), Equals, video.Size())
}
