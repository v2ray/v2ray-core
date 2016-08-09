// +build !coverage

package scenarios

import (
	"os"
	"os/exec"
)

func BuildV2Ray() error {
	binaryPath := GetTestBinaryPath()
	if _, err := os.Stat(binaryPath); err == nil {
		return nil
	}

	cmd := exec.Command("go", "build", "-tags=json", "-o="+binaryPath, GetSourcePath())
	return cmd.Run()
}

func RunV2Ray(configFile string) *exec.Cmd {
	binaryPath := GetTestBinaryPath()
	proc := exec.Command(binaryPath, "-config", configFile)
	proc.Stderr = os.Stderr
	proc.Stdout = os.Stdout

	return proc
}
