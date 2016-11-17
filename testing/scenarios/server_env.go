package scenarios

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"v2ray.com/core/common/log"

	"fmt"
	"io/ioutil"
	"sync"
	_ "v2ray.com/core"
	"v2ray.com/core/common/retry"
)

var (
	runningServers    = make([]*exec.Cmd, 0, 10)
	testBinaryPath    string
	testBinaryPathGen sync.Once
)

func GenTestBinaryPath() {
	testBinaryPathGen.Do(func() {
		var tempDir string
		err := retry.Timed(5, 100).On(func() error {
			dir, err := ioutil.TempDir("", "v2ray")
			if err != nil {
				return err
			}
			tempDir = dir
			return nil
		})
		if err != nil {
			panic(err)
		}
		file := filepath.Join(tempDir, "v2ray.test")
		if runtime.GOOS == "windows" {
			file += ".exe"
		}
		testBinaryPath = file
		fmt.Printf("Generated binary path: %s\n", file)
	})
}

func GetSourcePath() string {
	return filepath.Join("v2ray.com", "core", "main")
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
