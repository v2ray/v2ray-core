package main

import (
	"flag"
	"fmt"
	"os"

	commlog "v2ray.com/core/common/log"
	// _ "v2ray.com/core/infra/conf/command"
	"v2ray.com/core/infra/control"
)

func getCommandName() string {
	if len(os.Args) > 1 {
		return os.Args[1]
	}
	return ""
}

func main() {
	// let the v2ctl prints log at stderr
	commlog.RegisterHandler(commlog.NewLogger(commlog.CreateStderrLogWriter()))
	name := getCommandName()
	cmd := control.GetCommand(name)
	if cmd == nil {
		fmt.Fprintln(os.Stderr, "Unknown command:", name)
		fmt.Fprintln(os.Stderr)

		fmt.Println("v2ctl <command>")
		fmt.Println("Available commands:")
		control.PrintUsage()
		return
	}

	if err := cmd.Execute(os.Args[2:]); err != nil {
		hasError := false
		if err != flag.ErrHelp {
			fmt.Fprintln(os.Stderr, err.Error())
			fmt.Fprintln(os.Stderr)
			hasError = true
		}

		for _, line := range cmd.Description().Usage {
			fmt.Println(line)
		}

		if hasError {
			os.Exit(-1)
		}
	}
}
