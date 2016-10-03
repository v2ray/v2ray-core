package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"v2ray.com/core"
	_ "v2ray.com/core/app/router/rules"
	"v2ray.com/core/common/log"
	"v2ray.com/core/shell/point"

	// The following are necessary as they register handlers in their init functions.
	_ "v2ray.com/core/proxy/blackhole"
	_ "v2ray.com/core/proxy/dokodemo"
	_ "v2ray.com/core/proxy/freedom"
	_ "v2ray.com/core/proxy/http"
	_ "v2ray.com/core/proxy/shadowsocks"
	_ "v2ray.com/core/proxy/socks"
	_ "v2ray.com/core/proxy/vmess/inbound"
	_ "v2ray.com/core/proxy/vmess/outbound"

	_ "v2ray.com/core/transport/internet/kcp"
	_ "v2ray.com/core/transport/internet/tcp"
	_ "v2ray.com/core/transport/internet/udp"
	_ "v2ray.com/core/transport/internet/ws"

	_ "v2ray.com/core/transport/internet/authenticators/noop"
	_ "v2ray.com/core/transport/internet/authenticators/srtp"
	_ "v2ray.com/core/transport/internet/authenticators/utp"
)

var (
	configFile string
	logLevel   = flag.String("loglevel", "warning", "Level of log info to be printed to console, available value: debug, info, warning, error")
	version    = flag.Bool("version", false, "Show current version of V2Ray.")
	test       = flag.Bool("test", false, "Test config file only, without launching V2Ray server.")
	format     = flag.String("format", "json", "Format of input file.")
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
	var configInput io.Reader
	if configFile == "stdin:" {
		configInput = os.Stdin
	} else {
		fixedFile := os.ExpandEnv(configFile)
		file, err := os.Open(fixedFile)
		if err != nil {
			log.Error("Config file not readable: ", err)
			return nil
		}
		defer file.Close()
		configInput = file
	}
	config, err := point.LoadConfig(configInput)
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
		signal.Notify(osSignals, os.Interrupt, os.Kill, syscall.SIGTERM)

		<-osSignals
		point.Close()
	}
	log.Close()
}
