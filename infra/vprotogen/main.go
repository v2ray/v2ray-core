package main

import (
	"flag"
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

var (
	repo = flag.String("repo", "", "Repo for protobuf generation, such as v2ray.com/core")
)

func main() {
	flag.Parse()

	protofiles := make(map[string][]string)
	protoc := protocMap[runtime.GOOS]
	gosrc := filepath.Join(os.Getenv("GOPATH"), "src")
	reporoot := filepath.Join(os.Getenv("GOPATH"), "src", *repo)

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

	var protoFilesUsingProtocGenGoFast = map[string]bool{"proxy/vless/encoding/addons.proto": true}

	for _, files := range protofiles {
		for _, absPath := range files {
			relPath, _ := filepath.Rel(reporoot, absPath)
			args := make([]string, 0)
			if protoFilesUsingProtocGenGoFast[relPath] {
				args = []string{"--proto_path", reporoot, "--gofast_out", gosrc}
			} else {
				args = []string{"--proto_path", reporoot, "--go_out", gosrc, "--go-grpc_out", gosrc}
			}
			args = append(args, absPath)
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
}
