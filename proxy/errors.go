package proxy

import (
	"errors"
)

var (
	InvalidAuthentication  = errors.New("Invalid authentication.")
	InvalidProtocolVersion = errors.New("Invalid protocol version.")
	ErrorAlreadyListening  = errors.New("Already listening on another port.")
)
