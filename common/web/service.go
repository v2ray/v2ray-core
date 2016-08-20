package web

import (
	"v2ray.com/core/common/protocol"
)

type Authenciation struct {
	Required bool
	User     *protocol.User
}
