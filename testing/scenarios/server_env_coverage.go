// +build coverage

package scenarios

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/v2ray/v2ray-core/common/uuid"
)

func BuildV2Ray() error {
	binaryPath := GetTestBinaryPath()
	if _, err := os.Stat(binaryPath); err == nil {
		return nil
	}

	cmd := exec.Command("go", "test", "-tags", "json coverage coveragemain", "-coverpkg", "github.com/v2ray/v2ray-core/...", "-c", "-o", binaryPath, GetSourcePath())
	return cmd.Run()
}

func RunV2Ray(configFile string) *exec.Cmd {
	binaryPath := GetTestBinaryPath()

	covDir := filepath.Join(os.Getenv("GOPATH"), "out", "v2ray", "cov")
	profile := uuid.New().String() + ".out"
	proc := exec.Command(binaryPath, "-config", configFile, "-test.run", "TestRunMainForCoverage", "-test.coverprofile", profile, "-test.outputdir", covDir)
	proc.Stderr = os.Stderr
	proc.Stdout = os.Stdout

	return proc
}
