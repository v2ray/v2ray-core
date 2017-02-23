package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func copyConfigFile(src, dest string, goOS GoOS, format bool) error {
	content, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}
	if format {
		str := string(content)
		str = strings.Replace(str, "\r\n", "\n", -1)
		if goOS == Windows {
			str = strings.Replace(str, "\n", "\r\n", -1)
		}
		content = []byte(str)
	}
	return ioutil.WriteFile(dest, content, 0777)
}

func copyConfigFiles(dir string, goOS GoOS) error {
	GOPATH := os.Getenv("GOPATH")
	verifyDir := filepath.Join(GOPATH, "src", "v2ray.com", "core", "tools", "release", "verify")
	if err := copyConfigFile(filepath.Join(verifyDir, "official_release.asc"), filepath.Join(dir, "official_release.asc"), goOS, false); err != nil {
		return err
	}

	srcDir := filepath.Join(GOPATH, "src", "v2ray.com", "core", "tools", "release", "config")
	src := filepath.Join(srcDir, "vpoint_socks_vmess.json")
	dest := filepath.Join(dir, "vpoint_socks_vmess.json")
	if goOS == Windows || goOS == MacOS {
		dest = filepath.Join(dir, "config.json")
	}
	if err := copyConfigFile(src, dest, goOS, true); err != nil {
		return err
	}

	if goOS == Windows || goOS == MacOS {
		return nil
	}

	src = filepath.Join(srcDir, "vpoint_vmess_freedom.json")
	dest = filepath.Join(dir, "vpoint_vmess_freedom.json")

	if err := copyConfigFile(src, dest, goOS, true); err != nil {
		return err
	}

	if goOS == Linux {
		if err := os.MkdirAll(filepath.Join(dir, "systemv"), os.ModeDir|0777); err != nil {
			return err
		}
		src = filepath.Join(srcDir, "systemv", "v2ray")
		dest = filepath.Join(dir, "systemv", "v2ray")
		if err := copyConfigFile(src, dest, goOS, false); err != nil {
			return err
		}

		if err := os.MkdirAll(filepath.Join(dir, "systemd"), os.ModeDir|0777); err != nil {
			return err
		}
		src = filepath.Join(srcDir, "systemd", "v2ray.service")
		dest = filepath.Join(dir, "systemd", "v2ray.service")
		if err := copyConfigFile(src, dest, goOS, false); err != nil {
			return err
		}
	}

	return nil
}
