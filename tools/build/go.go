package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

func buildV2Ray(targetFile string, version string, goOS GoOS, goArch GoArch) error {
	ldFlags := "-s"
	if version != "custom" {
		year, month, day := time.Now().UTC().Date()
		today := fmt.Sprintf("%04d%02d%02d", year, int(month), day)
		ldFlags = ldFlags + " -X v2ray.com/core.version=" + version + " -X v2ray.com/core.build=" + today
	}
	cmd := exec.Command("go", "build", "-tags", "json", "-o", targetFile, "-compiler", "gc", "-ldflags", ldFlags, "v2ray.com/core/main")
	cmd.Env = append(cmd.Env, "GOOS="+string(goOS), "GOARCH="+string(goArch))
	cmd.Env = append(cmd.Env, os.Environ()...)
	output, err := cmd.CombinedOutput()
	if len(output) > 0 {
		fmt.Println(string(output))
	}
	return err
}
