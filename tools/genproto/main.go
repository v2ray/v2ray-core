// +build generate

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

var protocMap = map[string]string{
	"windows": filepath.Join(os.Getenv("GOPATH"), "src", "v2ray.com", "core", ".dev", "protoc", "windows", "protoc.exe"),
	"darwin":  filepath.Join(os.Getenv("GOPATH"), "src", "v2ray.com", "core", ".dev", "protoc", "macos", "protoc"),
	"linux":   filepath.Join(os.Getenv("GOPATH"), "src", "v2ray.com", "core", ".dev", "protoc", "linux", "protoc"),
}

func sdkPath(reporoot, lang string) string {
	path := filepath.Join(reporoot, ".dev", "sdk", lang)
	os.MkdirAll(path, os.ModePerm)
	return path
}

func main() {
	protofiles := make(map[string][]string)
	protoc := protocMap[runtime.GOOS]
	gosrc := filepath.Join(os.Getenv("GOPATH"), "src")
	reporoot := filepath.Join(gosrc, "v2ray.com", "core")

	filepath.Walk(reporoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}

		if info.IsDir() {
			return nil
		}

		dir := filepath.Dir(path)
		filename := filepath.Base(path)
		if strings.HasSuffix(filename, ".proto") {
			protofiles[dir] = append(protofiles[dir], path)
		}

		return nil
	})

	for _, files := range protofiles {
		args := []string{"--proto_path", gosrc, "--go_out", gosrc}
		args = append(args, files...)
		cmd := exec.Command(protoc, args...)
		cmd.Env = append(cmd.Env, os.Environ()...)
		output, err := cmd.CombinedOutput()
		if len(output) > 0 {
			fmt.Println(string(output))
		}
		if err != nil {
			fmt.Println(err)
		}
	}

	err := filepath.Walk(reporoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}

		if info.IsDir() {
			return nil
		}

		if !strings.HasSuffix(info.Name(), ".pb.go") {
			return nil
		}

		content, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		pos := bytes.Index(content, []byte("\npackage"))
		if pos > 0 {
			if err := ioutil.WriteFile(path, content[pos+1:], info.Mode()); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		panic(err)
	}
}
