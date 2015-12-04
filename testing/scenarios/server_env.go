package scenarios

import (
	"os"
	"path/filepath"

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
	serverup = make(map[string]bool)
)

func TestFile(filename string) string {
	return filepath.Join(os.Getenv("GOPATH"), "src", "github.com", "v2ray", "v2ray-core", "testing", "scenarios", "data", filename)
}

func InitializeServerSetOnce(testcase string) error {
	if up, found := serverup[testcase]; found && up {
		return nil
	}
	err := InitializeServer(TestFile(testcase + "_server.json"))
	if err != nil {
		return err
	}
	err = InitializeServer(TestFile(testcase + "_client.json"))
	if err != nil {
		return err
	}
	serverup[testcase] = true
	return nil
}

func InitializeServer(configFile string) error {
	config, err := jsonconf.LoadConfig(configFile)
	if err != nil {
		log.Error("Failed to read config file (%s): %v", configFile, err)
		return err
	}

	vPoint, err := point.NewPoint(config)
	if err != nil {
		log.Error("Failed to create Point server: %v", err)
		return err
	}

	err = vPoint.Start()
	if err != nil {
		log.Error("Error starting Point server: %v", err)
		return err
	}

	return nil
}
