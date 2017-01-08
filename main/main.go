package main

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
	"v2ray.com/core/common/log"

	_ "v2ray.com/core/main/distro/all"
	conf "v2ray.com/core/tools/conf"
)

var (
	configFile string
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

func startV2Ray() *core.Point {
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
	config, err := core.LoadConfig(GetConfigFormat(), configInput)
	if err != nil {
		log.Error("Failed to read config file (", configFile, "): ", configFile, err)
		return nil
	}

	vPoint, err := core.NewPoint(config)
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

	//开发期间自动引入json库（其实代码里面，什么都没有写）
	conf.ImportJsonParser()

	if point := startV2Ray(); point != nil {
		osSignals := make(chan os.Signal, 1)
		signal.Notify(osSignals, os.Interrupt, os.Kill, syscall.SIGTERM)

		<-osSignals
		point.Close()
	}
	log.Close()
}
