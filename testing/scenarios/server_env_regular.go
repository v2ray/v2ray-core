// +build !coverage

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

	cmd := exec.Command("go", "build", "-tags=json", "-o="+binaryPath, GetSourcePath())
	return cmd.Run()
}

func RunV2Ray(configFile string) *exec.Cmd {
	proc := exec.Command(binaryPath, "-config="+configFile)
	proc.Stderr = os.Stderr
	proc.Stdout = os.Stdout

	return proc
}
