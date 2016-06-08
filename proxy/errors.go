package proxy

import (
	"errors"
)

var (
	ErrorInvalidAuthentication  = errors.New("Invalid authentication.")
	ErrorInvalidProtocolVersion = errors.New("Invalid protocol version.")
	ErrorAlreadyListening       = errors.New("Already listening on another port.")
)
