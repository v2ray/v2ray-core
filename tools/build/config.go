package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func copyConfigFile(src, dest string, goOS GoOS) error {
	content, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}
	str := string(content)
	str = strings.Replace(str, "\r\n", "\n", -1)
	if goOS == Windows {
		str = strings.Replace(str, "\n", "\r\n", -1)
	}
	return ioutil.WriteFile(dest, []byte(str), 0777)
}

func copyConfigFiles(dir string, goOS GoOS) error {
	GOPATH := os.Getenv("GOPATH")
	srcDir := filepath.Join(GOPATH, "src", "github.com", "v2ray", "v2ray-core", "release", "config")
	src := filepath.Join(srcDir, "vpoint_socks_vmess.json")
	dest := filepath.Join(dir, "vpoint_socks_vmess.json")
	if goOS == Windows || goOS == MacOS {
		dest = filepath.Join(dir, "config.json")
	}
	err := copyConfigFile(src, dest, goOS)
	if err != nil {
		return err
	}

	if goOS == Windows || goOS == MacOS {
		return nil
	}

	src = filepath.Join(srcDir, "vpoint_vmess_freedom.json")
	dest = filepath.Join(dir, "vpoint_vmess_freedom.json")
	return copyConfigFile(src, dest, goOS)
}
