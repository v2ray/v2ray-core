package protocol

import (
	"errors"
)

var (
	ErrInvalidUser    = errors.New("Invalid user.")
	ErrInvalidVersion = errors.New("Invalid version.")
)
