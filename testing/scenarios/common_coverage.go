// +build coverage

package scenarios

import (
	"bytes"
	"os"
	"os/exec"

	"v2ray.com/core/common/uuid"
)

func BuildV2Ray() error {
	genTestBinaryPath()
	if _, err := os.Stat(testBinaryPath); err == nil {
		return nil
	}

	cmd := exec.Command("go", "test", "-tags", "coverage coveragemain", "-coverpkg", "v2ray.com/core/...", "-c", "-o", testBinaryPath, GetSourcePath())
	return cmd.Run()
}

func RunV2RayProtobuf(config []byte) *exec.Cmd {
	genTestBinaryPath()

	covDir := os.Getenv("V2RAY_COV")
	os.MkdirAll(covDir, os.ModeDir)
	randomID := uuid.New()
	profile := randomID.String() + ".out"
	proc := exec.Command(testBinaryPath, "-config=stdin:", "-format=pb", "-test.run", "TestRunMainForCoverage", "-test.coverprofile", profile, "-test.outputdir", covDir)
	proc.Stdin = bytes.NewBuffer(config)
	proc.Stderr = os.Stderr
	proc.Stdout = os.Stdout

	return proc
}
