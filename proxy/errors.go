package proxy

import (
	"v2ray.com/core/common/errors"
)

var (
	ErrInvalidProtocolVersion = errors.New("Invalid protocol version.")
	ErrAlreadyListening       = errors.New("Already listening on another port.")
)
