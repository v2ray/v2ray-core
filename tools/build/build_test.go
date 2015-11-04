package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/v2ray/v2ray-core/testing/unit"
)

func cleanBinPath() {
	os.RemoveAll(binPath)
	os.Mkdir(binPath, os.ModeDir|0777)
}

func fileExists(file string) bool {
	_, err := os.Stat(file)
	return err == nil
}

func allFilesExists(files ...string) bool {
	for _, file := range files {
		fullPath := filepath.Join(binPath, file)
		if !fileExists(fullPath) {
			fmt.Println(fullPath + " doesn't exist.")
			return false
		}
	}
	return true
}

func TestBuildMacOS(t *testing.T) {
	assert := unit.Assert(t)
	binPath = filepath.Join(os.Getenv("GOPATH"), "testing")
	cleanBinPath()

	build("macos", "amd64", true, "test")
	assert.Bool(allFilesExists(
		"v2ray-macos.zip",
		"v2ray-test-macos",
		filepath.Join("v2ray-test-macos", "config.json"),
		filepath.Join("v2ray-test-macos", "v2ray"))).IsTrue()

	build("windows", "amd64", true, "test")
	assert.Bool(allFilesExists(
		"v2ray-windows-64.zip",
		"v2ray-test-windows-64",
		filepath.Join("v2ray-test-windows-64", "config.json"),
		filepath.Join("v2ray-test-windows-64", "v2ray.exe"))).IsTrue()

	build("linux", "amd64", true, "test")
	assert.Bool(allFilesExists(
		"v2ray-linux-64.zip",
		"v2ray-test-linux-64",
		filepath.Join("v2ray-test-linux-64", "vpoint_socks_vmess.json"),
		filepath.Join("v2ray-test-linux-64", "vpoint_vmess_freedom.json"),
		filepath.Join("v2ray-test-linux-64", "v2ray"))).IsTrue()
}
