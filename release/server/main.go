package main

import (
	"flag"
	"fmt"

	"github.com/v2ray/v2ray-core"
  jsonconf "github.com/v2ray/v2ray-core/io/config/json"
	"github.com/v2ray/v2ray-core/log"

	// The following are neccesary as they register handlers in their init functions.
	_ "github.com/v2ray/v2ray-core/net/freedom"
	_ "github.com/v2ray/v2ray-core/net/socks"
	_ "github.com/v2ray/v2ray-core/net/vmess"
)

var (
	configFile = flag.String("config", "", "Config file for this Point server.")
	logLevel   = flag.String("loglevel", "", "Level of log info to be printed to console, available value: debug, info, warning, error")
	version    = flag.Bool("version", false, "Show current version of V2Ray.")
)

func main() {
	flag.Parse()

	if *version {
		fmt.Printf("V2Ray version %s (%s): %s", core.Version, core.Codename, core.Intro)
		fmt.Println()
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
	}

	if configFile == nil || len(*configFile) == 0 {
		panic(log.Error("Config file is not set."))
	}
  config, err := jsonconf.LoadConfig(*configFile)
	if err != nil {
		panic(log.Error("Failed to read config file (%s): %v", *configFile, err))
	}

	vPoint, err := core.NewPoint(config)
	if err != nil {
		panic(log.Error("Failed to create Point server: %v", err))
	}

	err = vPoint.Start()
	if err != nil {
		log.Error("Error starting Point server: %v", err)
	}

	finish := make(chan bool)
	<-finish
}
