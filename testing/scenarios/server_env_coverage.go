// +build coverage

package scenarios

import (
	"os"
	"os/exec"
)

func BuildV2Ray() error {
	if _, err := os.Stat(binaryPath); err == nil {
		return nil
	}

	if err := FillBinaryPath(); err != nil {
		return err
	}

	cmd := exec.Command("go", "test", "-tags", "json coverage coveragemain", "-coverpkg", "github.com/v2ray/v2ray-core/...", "-c", "-o", binaryPath, GetSourcePath())
	return cmd.Run()
}

func RunV2Ray(configFile string) *exec.Cmd {
	proc := exec.Command(binaryPath, "-config="+configFile, "-test.run=TestRunMainForCoverage", "-test.coverprofile=coversingle.out")
	proc.Stderr = os.Stderr
	proc.Stdout = os.Stdout

	return proc
}
