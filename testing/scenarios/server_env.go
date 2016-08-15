package scenarios

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	_ "github.com/v2ray/v2ray-core/app/router/rules"
	"github.com/v2ray/v2ray-core/common/log"

	// The following are necessary as they register handlers in their init functions.
	_ "github.com/v2ray/v2ray-core/proxy/blackhole"
	_ "github.com/v2ray/v2ray-core/proxy/dokodemo"
	_ "github.com/v2ray/v2ray-core/proxy/freedom"
	_ "github.com/v2ray/v2ray-core/proxy/http"
	_ "github.com/v2ray/v2ray-core/proxy/shadowsocks"
	_ "github.com/v2ray/v2ray-core/proxy/socks"
	_ "github.com/v2ray/v2ray-core/proxy/vmess/inbound"
	_ "github.com/v2ray/v2ray-core/proxy/vmess/outbound"
)

var (
	runningServers = make([]*exec.Cmd, 0, 10)
)

func GetTestBinaryPath() string {
	file := filepath.Join(os.Getenv("GOPATH"), "out", "v2ray", "v2ray.test")
	if runtime.GOOS == "windows" {
		file += ".exe"
	}
	return file
}

func GetSourcePath() string {
	return filepath.Join("github.com", "v2ray", "v2ray-core", "shell", "point", "main")
}

func TestFile(filename string) string {
	return filepath.Join(os.Getenv("GOPATH"), "src", "github.com", "v2ray", "v2ray-core", "testing", "scenarios", "data", filename)
}

func InitializeServerSetOnce(testcase string) error {
	if err := InitializeServerServer(testcase); err != nil {
		return err
	}
	if err := InitializeServerClient(testcase); err != nil {
		return err
	}
	return nil
}

func InitializeServerServer(testcase string) error {
	return InitializeServer(TestFile(testcase + "_server.json"))
}

func InitializeServerClient(testcase string) error {
	return InitializeServer(TestFile(testcase + "_client.json"))
}

func InitializeServer(configFile string) error {
	err := BuildV2Ray()
	if err != nil {
		return err
	}

	proc := RunV2Ray(configFile)

	err = proc.Start()
	if err != nil {
		return err
	}

	time.Sleep(time.Second)

	runningServers = append(runningServers, proc)

	return nil
}

func CloseAllServers() {
	log.Info("Closing all servers.")
	for _, server := range runningServers {
		server.Process.Signal(os.Interrupt)
		server.Process.Wait()
	}
	runningServers = make([]*exec.Cmd, 0, 10)
	log.Info("All server closed.")
}
