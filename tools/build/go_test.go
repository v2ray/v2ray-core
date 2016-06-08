package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/v2ray/v2ray-core/testing/assert"
)

func TestBuildAndRun(t *testing.T) {
	assert := assert.On(t)

	gopath := os.Getenv("GOPATH")
	goOS := parseOS(runtime.GOOS)
	goArch := parseArch(runtime.GOARCH)
	target := filepath.Join(gopath, "src", "v2ray_test")
	if goOS == Windows {
		target += ".exe"
	}
	err := buildV2Ray(target, "v1.0", goOS, goArch)
	assert.Error(err).IsNil()

	outBuffer := bytes.NewBuffer(make([]byte, 0, 1024))
	errBuffer := bytes.NewBuffer(make([]byte, 0, 1024))
	configFile := filepath.Join(gopath, "src", "github.com", "v2ray", "v2ray-core", "release", "config", "vpoint_socks_vmess.json")
	cmd := exec.Command(target, "--config="+configFile)
	cmd.Stdout = outBuffer
	cmd.Stderr = errBuffer
	cmd.Start()

	<-time.After(1 * time.Second)
	cmd.Process.Kill()

	outStr := string(outBuffer.Bytes())
	errStr := string(errBuffer.Bytes())

	assert.Bool(strings.Contains(outStr, "v1.0")).IsTrue()
	assert.String(errStr).Equals("")

	os.Remove(target)
}
