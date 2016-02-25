package web

import (
	"github.com/v2ray/v2ray-core/common/protocol"
)

type Authenciation struct {
	Required bool
	User     *protocol.User
}
