package command

import (
	"errors"
	"io"
)

var (
	ErrorNoSuchCommand = errors.New("No such command.")
)

type Command interface {
	Marshal(io.Writer)
	Unmarshal([]byte) error
}

type CommandCreator func() Command
