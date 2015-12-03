package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/v2ray/v2ray-core"
	_ "github.com/v2ray/v2ray-core/app/router/config/json"
	_ "github.com/v2ray/v2ray-core/app/router/rules"
	_ "github.com/v2ray/v2ray-core/app/router/rules/config/json"
	"github.com/v2ray/v2ray-core/common/log"
	"github.com/v2ray/v2ray-core/shell/point"
	jsonconf "github.com/v2ray/v2ray-core/shell/point/config/json"

	// The following are neccesary as they register handlers in their init functions.
	_ "github.com/v2ray/v2ray-core/proxy/blackhole"
	_ "github.com/v2ray/v2ray-core/proxy/blackhole/config/json"
	_ "github.com/v2ray/v2ray-core/proxy/dokodemo"
	_ "github.com/v2ray/v2ray-core/proxy/dokodemo/config/json"
	_ "github.com/v2ray/v2ray-core/proxy/freedom"
	_ "github.com/v2ray/v2ray-core/proxy/freedom/config/json"
	_ "github.com/v2ray/v2ray-core/proxy/socks"
	_ "github.com/v2ray/v2ray-core/proxy/socks/config/json"
	_ "github.com/v2ray/v2ray-core/proxy/vmess"
	_ "github.com/v2ray/v2ray-core/proxy/vmess/config/json"
)

var (
	configFile string
	logLevel   = flag.String("loglevel", "warning", "Level of log info to be printed to console, available value: debug, info, warning, error")
	version    = flag.Bool("version", false, "Show current version of V2Ray.")
)

func init() {
	defaultConfigFile := ""
	workingDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err == nil {
		defaultConfigFile = filepath.Join(workingDir, "config.json")
	}
	flag.StringVar(&configFile, "config", defaultConfigFile, "Config file for this Point server.")
}

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

	if len(configFile) == 0 {
		log.Error("Config file is not set.")
		return
	}
	config, err := jsonconf.LoadConfig(configFile)
	if err != nil {
		log.Error("Failed to read config file (%s): %v", configFile, err)
		return
	}

	if config.LogConfig() != nil && len(config.LogConfig().AccessLog()) > 0 {
		log.InitAccessLogger(config.LogConfig().AccessLog())
	}

	vPoint, err := point.NewPoint(config)
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
