// +build !coverage

package scenarios

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

func BuildV2Ray() error {
	GenTestBinaryPath()
	if _, err := os.Stat(testBinaryPath); err == nil {
		return nil
	}

	fmt.Printf("Building V2Ray into path (%s)\n", testBinaryPath)
	cmd := exec.Command("go", "build", "-tags=json", "-o="+testBinaryPath, GetSourcePath())
	return cmd.Run()
}

func RunV2Ray(configFile string) *exec.Cmd {
	GenTestBinaryPath()
	proc := exec.Command(testBinaryPath, "-config", configFile)
	proc.Stderr = os.Stderr
	proc.Stdout = os.Stdout

	return proc
}

func RunV2RayProtobuf(config []byte) *exec.Cmd {
	GenTestBinaryPath()
	proc := exec.Command(testBinaryPath, "-config=stdin:", "-format=pb")
	proc.Stdin = bytes.NewBuffer(config)
	proc.Stderr = os.Stderr
	proc.Stdout = os.Stdout

	return proc
}
