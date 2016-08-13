package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/v2ray/v2ray-core"
	_ "github.com/v2ray/v2ray-core/app/router/rules"
	"github.com/v2ray/v2ray-core/common/log"
	"github.com/v2ray/v2ray-core/shell/point"

	// The following are necessary as they register handlers in their init functions.
	_ "github.com/v2ray/v2ray-core/proxy/blackhole"
	_ "github.com/v2ray/v2ray-core/proxy/dokodemo"
	_ "github.com/v2ray/v2ray-core/proxy/freedom"
	_ "github.com/v2ray/v2ray-core/proxy/http"
	_ "github.com/v2ray/v2ray-core/proxy/shadowsocks"
	_ "github.com/v2ray/v2ray-core/proxy/socks"
	_ "github.com/v2ray/v2ray-core/proxy/vmess/inbound"
	_ "github.com/v2ray/v2ray-core/proxy/vmess/outbound"

	_ "github.com/v2ray/v2ray-core/transport/internet/kcp"
	_ "github.com/v2ray/v2ray-core/transport/internet/tcp"
	_ "github.com/v2ray/v2ray-core/transport/internet/udp"
	_ "github.com/v2ray/v2ray-core/transport/internet/ws"

	_ "github.com/v2ray/v2ray-core/transport/internet/authenticators/noop"
	_ "github.com/v2ray/v2ray-core/transport/internet/authenticators/srtp"
)

var (
	configFile string
	logLevel   = flag.String("loglevel", "warning", "Level of log info to be printed to console, available value: debug, info, warning, error")
	version    = flag.Bool("version", false, "Show current version of V2Ray.")
	test       = flag.Bool("test", false, "Test config file only, without launching V2Ray server.")
)

func init() {
	defaultConfigFile := ""
	workingDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err == nil {
		defaultConfigFile = filepath.Join(workingDir, "config.json")
	}
	flag.StringVar(&configFile, "config", defaultConfigFile, "Config file for this Point server.")
}

func startV2Ray() *point.Point {
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
		return nil
	}

	if len(configFile) == 0 {
		log.Error("Config file is not set.")
		return nil
	}
	config, err := point.LoadConfig(configFile)
	if err != nil {
		log.Error("Failed to read config file (", configFile, "): ", configFile, err)
		return nil
	}

	if config.LogConfig != nil && len(config.LogConfig.AccessLog) > 0 {
		log.InitAccessLogger(config.LogConfig.AccessLog)
	}

	vPoint, err := point.NewPoint(config)
	if err != nil {
		log.Error("Failed to create Point server: ", err)
		return nil
	}

	if *test {
		fmt.Println("Configuration OK.")
		return nil
	}

	err = vPoint.Start()
	if err != nil {
		log.Error("Error starting Point server: ", err)
		return nil
	}

	return vPoint
}

func main() {
	flag.Parse()

	core.PrintVersion()

	if *version {
		return
	}

	if point := startV2Ray(); point != nil {
		osSignals := make(chan os.Signal, 1)
		signal.Notify(osSignals, os.Interrupt, os.Kill)

		<-osSignals
		point.Close()
	}
	log.Close()
}
