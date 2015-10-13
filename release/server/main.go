package main

import (
	"flag"
	"fmt"

	"github.com/v2ray/v2ray-core"
	"github.com/v2ray/v2ray-core/common/log"
	jsonconf "github.com/v2ray/v2ray-core/config/json"

	// The following are neccesary as they register handlers in their init functions.
	_ "github.com/v2ray/v2ray-core/proxy/freedom"
	_ "github.com/v2ray/v2ray-core/proxy/freedom/config/json"
	_ "github.com/v2ray/v2ray-core/proxy/socks"
	_ "github.com/v2ray/v2ray-core/proxy/vmess"
)

var (
	configFile = flag.String("config", "", "Config file for this Point server.")
	logLevel   = flag.String("loglevel", "warning", "Level of log info to be printed to console, available value: debug, info, warning, error")
	version    = flag.Bool("version", false, "Show current version of V2Ray.")
)

func main() {
	flag.Parse()

	core.PrintVersion()

	if *version {
		return
	}

	switch *logLevel {
	case "debug":
		log.SetLogLevel(log.DebugLevel)
	case "info":
		log.SetLogLevel(log.InfoLevel)
	case "warning":
		log.SetLogLevel(log.WarningLevel)
	case "error":
		log.SetLogLevel(log.ErrorLevel)
	default:
		fmt.Println("Unknown log level: " + *logLevel)
		return
	}

	if configFile == nil || len(*configFile) == 0 {
		log.Error("Config file is not set.")
		return
	}
	config, err := jsonconf.LoadConfig(*configFile)
	if err != nil {
		log.Error("Failed to read config file (%s): %v", *configFile, err)
		return
	}

	if config.LogConfig() != nil && len(config.LogConfig().AccessLog()) > 0 {
		log.InitAccessLogger(config.LogConfig().AccessLog())
	}

	vPoint, err := core.NewPoint(config)
	if err != nil {
		log.Error("Failed to create Point server: %v", err)
		return
	}

	err = vPoint.Start()
	if err != nil {
		log.Error("Error starting Point server: %v", err)
		return
	}

	finish := make(chan bool)
	<-finish
}
