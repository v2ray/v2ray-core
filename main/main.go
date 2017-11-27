package main

//go:generate go run $GOPATH/src/v2ray.com/core/tools/generrorgen/main.go -pkg main -path Main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"v2ray.com/core"

	_ "v2ray.com/core/main/distro/all"
)

var (
	configFile string
	version    = flag.Bool("version", false, "Show current version of V2Ray.")
	test       = flag.Bool("test", false, "Test config file only, without launching V2Ray server.")
	format     = flag.String("format", "json", "Format of input file.")
	plugin     = flag.Bool("plugin", false, "True to load plugins.")
)

func init() {
	defaultConfigFile := ""
	workingDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err == nil {
		defaultConfigFile = filepath.Join(workingDir, "config.json")
	}
	flag.StringVar(&configFile, "config", defaultConfigFile, "Config file for this Point server.")
}

func GetConfigFormat() core.ConfigFormat {
	switch strings.ToLower(*format) {
	case "json":
		return core.ConfigFormat_JSON
	case "pb", "protobuf":
		return core.ConfigFormat_Protobuf
	default:
		return core.ConfigFormat_JSON
	}
}

func startV2Ray() (core.Server, error) {
	if len(configFile) == 0 {
		return nil, newError("config file is not set")
	}
	var configInput io.Reader
	if configFile == "stdin:" {
		configInput = os.Stdin
	} else {
		fixedFile := os.ExpandEnv(configFile)
		file, err := os.Open(fixedFile)
		if err != nil {
			return nil, newError("config file not readable").Base(err)
		}
		defer file.Close()
		configInput = file
	}
	config, err := core.LoadConfig(GetConfigFormat(), configInput)
	if err != nil {
		return nil, newError("failed to read config file: ", configFile).Base(err)
	}

	server, err := core.New(config)
	if err != nil {
		return nil, newError("failed to create server").Base(err)
	}

	return server, nil
}

func main() {
	flag.Parse()

	core.PrintVersion()

	if *version {
		return
	}

	if *plugin {
		if err := core.LoadPlugins(); err != nil {
			fmt.Println("Failed to load plugins:", err.Error())
			os.Exit(-1)
		}
	}

	server, err := startV2Ray()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}

	if *test {
		fmt.Println("Configuration OK.")
		os.Exit(0)
	}

	if err := server.Start(); err != nil {
		fmt.Println("Failed to start", err)
		os.Exit(-1)
	}

	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, os.Kill, syscall.SIGTERM)

	<-osSignals
	server.Close()
}
