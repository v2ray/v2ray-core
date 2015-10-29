package errors

import (
	"errors"
)

var (
	InvalidAuthentication  = errors.New("Invalid authentication.")
	InvalidProtocolVersion = errors.New("Invalid protocol version.")
)
