package main

//go:generate errorgen

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"

	"v2ray.com/core"
	"v2ray.com/core/common/cmdarg"
	"v2ray.com/core/common/platform"
	_ "v2ray.com/core/main/distro/all"
)

var (
	configFiles cmdarg.Arg // "Config file for V2Ray.", the option is customed type, parse in main
	configDir   string
	version     = flag.Bool("version", false, "Show current version of V2Ray.")
	test        = flag.Bool("test", false, "Test config file only, without launching V2Ray server.")
	format      = flag.String("format", "json", "Format of input file.")
	errNoConfig = newError("no valid config")
)

func fileExists(file string) bool {
	info, err := os.Stat(file)
	return err == nil && !info.IsDir()
}

func dirExists(file string) bool {
	info, err := os.Stat(file)
	return err == nil && info.IsDir()
}

func readConfDir(dirPath string) {
	confs, err := ioutil.ReadDir(dirPath)
	if err != nil {
		log.Fatalln(err)
	}
	for _, f := range confs {
		if strings.HasSuffix(f.Name(), ".json") {
			configFiles.Set(path.Join(dirPath, f.Name()))
		}
	}
}

func getConfigFilePath() (cmdarg.Arg, error) {
	if dirExists(configDir) {
		readConfDir(configDir)
	}
	if len(configFiles) > 0 {
		return configFiles, nil
	}

	if workingDir, err := os.Getwd(); err == nil {
		configFile := filepath.Join(workingDir, "config.json")
		if fileExists(configFile) {
			log.Println("Using default config: ", configFile)
			return cmdarg.Arg{configFile}, nil
		}
	}

	if configFile := platform.GetConfigurationPath(); fileExists(configFile) {
		log.Println("Using config from env: ", configFile)
		return cmdarg.Arg{configFile}, nil
	}

	if envConfDir := platform.GetConfDirPath(); dirExists(envConfDir) {
		log.Println("Using confdir from env: ", envConfDir)
		readConfDir(envConfDir)
		if len(configFiles) > 0 {
			return configFiles, nil
		}
	}

	log.Println("Using config from STDIN")
	return cmdarg.Arg{"stdin:"}, nil
}

func GetConfigFormat() string {
	switch strings.ToLower(*format) {
	case "pb", "protobuf":
		return "protobuf"
	default:
		return "json"
	}
}

func startV2Ray() (core.Server, error) {
	configFiles, err := getConfigFilePath()
	if err != nil {
		return nil, err
	}

	config, err := core.LoadConfig(GetConfigFormat(), configFiles[0], configFiles)
	if err != nil {
		return nil, newError("failed to read config files: [", configFiles.String(), "]").Base(err)
	}

	server, err := core.New(config)
	if err != nil {
		return nil, newError("failed to create server").Base(err)
	}

	return server, nil
}

func printVersion() {
	version := core.VersionStatement()
	for _, s := range version {
		fmt.Println(s)
	}
}

func main() {
	flag.Var(&configFiles, "config", "Config file for V2Ray. Multiple assign is accepted (only json). Latter ones overrides the former ones.")
	flag.Var(&configFiles, "c", "Short alias of -config")
	flag.StringVar(&configDir, "confdir", "", "A dir with multiple json config")
	flag.Parse()

	printVersion()

	if *version {
		return
	}

	server, err := startV2Ray()
	if err != nil {
		fmt.Println(err.Error())
		// Configuration error. Exit with a special value to prevent systemd from restarting.
		if err == errNoConfig {
			flag.PrintDefaults()
		}
		os.Exit(23)
	}

	if *test {
		fmt.Println("Configuration OK.")
		os.Exit(0)
	}

	if err := server.Start(); err != nil {
		fmt.Println("Failed to start", err)
		os.Exit(-1)
	}
	defer server.Close()

	// Explicitly triggering GC to remove garbage from config loading.
	runtime.GC()

	{
		osSignals := make(chan os.Signal, 1)
		signal.Notify(osSignals, os.Interrupt, os.Kill, syscall.SIGTERM)
		<-osSignals
	}
}
