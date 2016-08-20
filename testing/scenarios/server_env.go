package scenarios

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	_ "v2ray.com/core/app/router/rules"
	"v2ray.com/core/common/log"

	// The following are necessary as they register handlers in their init functions.
	_ "v2ray.com/core/proxy/blackhole"
	_ "v2ray.com/core/proxy/dokodemo"
	_ "v2ray.com/core/proxy/freedom"
	_ "v2ray.com/core/proxy/http"
	_ "v2ray.com/core/proxy/shadowsocks"
	_ "v2ray.com/core/proxy/socks"
	_ "v2ray.com/core/proxy/vmess/inbound"
	_ "v2ray.com/core/proxy/vmess/outbound"
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
	return filepath.Join("v2ray.com", "core", "shell", "point", "main")
}

func TestFile(filename string) string {
	return filepath.Join(os.Getenv("GOPATH"), "src", "v2ray.com", "core", "testing", "scenarios", "data", filename)
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
