package protocol

import (
	"errors"
)

var (
	ErrorInvalidUser    = errors.New("Invalid user.")
	ErrorInvalidVersion = errors.New("Invalid version.")
)
