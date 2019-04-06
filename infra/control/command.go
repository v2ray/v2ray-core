package control

import (
	"fmt"
	"strings"
)

type Description struct {
	Short string
	Usage []string
}

type Command interface {
	Name() string
	Description() Description
	Execute(args []string) error
}

var (
	commandRegistry = make(map[string]Command)
)

func RegisterCommand(cmd Command) error {
	entry := strings.ToLower(cmd.Name())
	if len(entry) == 0 {
		return newError("empty command name")
	}
	commandRegistry[entry] = cmd
	return nil
}

func GetCommand(name string) Command {
	cmd, found := commandRegistry[name]
	if !found {
		return nil
	}
	return cmd
}

type hiddenCommand interface {
	Hidden() bool
}

func PrintUsage() {
	for name, cmd := range commandRegistry {
		if _, ok := cmd.(hiddenCommand); ok {
			continue
		}
		fmt.Println("   ", name, "\t\t\t", cmd.Description())
	}
}
