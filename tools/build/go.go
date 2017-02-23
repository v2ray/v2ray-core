package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

func buildV2Ray(targetFile string, version string, goOS GoOS, goArch GoArch) error {
	goPath := os.Getenv("GOPATH")
	ldFlags := "-s"
	if version != "custom" {
		year, month, day := time.Now().UTC().Date()
		today := fmt.Sprintf("%04d%02d%02d", year, int(month), day)
		ldFlags = ldFlags + " -X v2ray.com/core.version=" + version + " -X v2ray.com/core.build=" + today
	}
	cmd := exec.Command(
		"go", "build",
		"-tags", "json",
		"-o", targetFile,
		"-compiler", "gc",
		"-ldflags", ldFlags,
		"-gcflags", "-trimpath="+goPath,
		"-asmflags", "-trimpath="+goPath,
		"v2ray.com/core/main")
	cmd.Env = append(cmd.Env, "GOOS="+string(goOS), "GOARCH="+string(goArch), "CGO_ENABLED=0")
	cmd.Env = append(cmd.Env, os.Environ()...)
	output, err := cmd.CombinedOutput()
	if len(output) > 0 {
		fmt.Println(string(output))
	}
	return err
}

func signFile(file string) error {
	pass := os.Getenv("GPG_SIGN_PASS")
	cmd := exec.Command("gpg", "--digest-algo", "SHA512", "--no-tty", "--batch", "--passphrase", pass, "--output", file+".sig", "--detach-sig", file)
	cmd.Env = append(cmd.Env, os.Environ()...)
	output, err := cmd.CombinedOutput()
	if len(output) > 0 {
		fmt.Println(string(output))
	}
	return err
}
