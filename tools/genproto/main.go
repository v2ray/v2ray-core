// +build generate

package main

import (
	"fmt"
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

func main() {
	protofiles := make(map[string][]string)
	protoc := protocMap[runtime.GOOS]
	gosrc := filepath.Join(os.Getenv("GOPATH"), "src")

	filepath.Walk(filepath.Join(gosrc, "v2ray.com", "core"), func(path string, info os.FileInfo, err error) error {
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
}
