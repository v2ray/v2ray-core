package main

import (
	"flag"
	"io/ioutil"
	"path"

	"github.com/v2ray/v2ray-core"
	"github.com/v2ray/v2ray-core/log"

	_ "github.com/v2ray/v2ray-core/net/freedom"
	_ "github.com/v2ray/v2ray-core/net/socks"
	_ "github.com/v2ray/v2ray-core/net/vmess"
)

var (
	configFile = flag.String("config", "", "Config file for this VPoint server.")
)

func main() {
	flag.Parse()

	log.SetLogLevel(log.DebugLevel)

	if configFile == nil || len(*configFile) == 0 {
		panic(log.Error("Config file is not set."))
	}
	rawVConfig, err := ioutil.ReadFile(*configFile)
	if err != nil {
		panic(log.Error("Failed to read config file (%s): %v", *configFile, err))
	}
	vconfig, err := core.LoadVConfig(rawVConfig)
	if err != nil {
		panic(log.Error("Failed to parse VConfig: %v", err))
	}

	if !path.IsAbs(vconfig.InboundConfig.File) && len(vconfig.InboundConfig.File) > 0 {
		vconfig.InboundConfig.File = path.Join(path.Dir(*configFile), vconfig.InboundConfig.File)
	}

	if !path.IsAbs(vconfig.OutboundConfig.File) && len(vconfig.OutboundConfig.File) > 0 {
		vconfig.OutboundConfig.File = path.Join(path.Dir(*configFile), vconfig.OutboundConfig.File)
	}

	vPoint, err := core.NewVPoint(vconfig)
	if err != nil {
		panic(log.Error("Failed to create VPoint server: %v", err))
	}

	err = vPoint.Start()
	if err != nil {
		log.Error("Error starting VPoint server.")
	}

	finish := make(chan bool)
	<-finish
}
