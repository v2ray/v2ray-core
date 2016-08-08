// +build coverage

package scenarios

import (
	"os"
	"os/exec"
	"path/filepath"
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
	profile := "coversingle.out"
	wd, err := os.Getwd()
	if err == nil {
		profile = filepath.Join(wd, profile)
	}
	proc := exec.Command(binaryPath, "-config", configFile, "-test.run", "TestRunMainForCoverage", "-test.coverprofile", profile)
	proc.Stderr = os.Stderr
	proc.Stdout = os.Stdout

	return proc
}
