package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/v2ray/v2ray-core"
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
	rawVConfig, err := ioutil.ReadFile(*configFile)
	if err != nil {
		panic(log.Error("Failed to read config file (%s): %v", *configFile, err))
	}
	vconfig, err := core.LoadConfig(rawVConfig)
	if err != nil {
		panic(log.Error("Failed to parse Config: %v", err))
	}

	if !filepath.IsAbs(vconfig.InboundConfig.File) && len(vconfig.InboundConfig.File) > 0 {
		vconfig.InboundConfig.File = filepath.Join(filepath.Dir(*configFile), vconfig.InboundConfig.File)
	}

	if !filepath.IsAbs(vconfig.OutboundConfig.File) && len(vconfig.OutboundConfig.File) > 0 {
		vconfig.OutboundConfig.File = filepath.Join(filepath.Dir(*configFile), vconfig.OutboundConfig.File)
	}

	vPoint, err := core.NewPoint(vconfig)
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
