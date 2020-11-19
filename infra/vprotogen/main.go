package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"v2ray.com/core"
	"v2ray.com/core/common"
)

func main() {
	pwd, wdErr := os.Getwd()
	if wdErr != nil {
		fmt.Println("Can not get current working directory.")
		os.Exit(1)
	}

	GOBIN := common.GetGOBIN()
	protoc := core.ProtocMap[runtime.GOOS]

	protoFilesMap := make(map[string][]string)
	walkErr := filepath.Walk("./", func(path string, info os.FileInfo, err error) error {
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
			protoFilesMap[dir] = append(protoFilesMap[dir], path)
		}

		return nil
	})
	if walkErr != nil {
		fmt.Println(walkErr)
		os.Exit(1)
	}

	for _, files := range protoFilesMap {
		for _, relProtoFile := range files {
			var args []string
			if core.ProtoFilesUsingProtocGenGoFast[relProtoFile] {
				args = []string{"--gofast_out", pwd, "--plugin", "protoc-gen-gofast=" + GOBIN + "/protoc-gen-gofast"}
			} else {
				args = []string{"--go_out", pwd, "--go-grpc_out", pwd, "--plugin", "protoc-gen-go=" + GOBIN + "/protoc-gen-go", "--plugin", "protoc-gen-go-grpc=" + GOBIN + "/protoc-gen-go-grpc"}
			}
			args = append(args, relProtoFile)
			cmd := exec.Command(protoc, args...)
			cmd.Env = append(cmd.Env, os.Environ()...)
			cmd.Env = append(cmd.Env, "GOBIN="+GOBIN)
			output, cmdErr := cmd.CombinedOutput()
			if len(output) > 0 {
				fmt.Println(string(output))
			}
			if cmdErr != nil {
				fmt.Println(cmdErr)
				os.Exit(1)
			}
		}
	}

	moduleName, gmnErr := common.GetModuleName(pwd)
	if gmnErr != nil {
		fmt.Println(gmnErr)
		os.Exit(1)
	}
	modulePath := filepath.Join(strings.Split(moduleName, "/")...)

	pbGoFilesMap := make(map[string][]string)
	walkErr2 := filepath.Walk(modulePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}

		if info.IsDir() {
			return nil
		}

		dir := filepath.Dir(path)
		filename := filepath.Base(path)
		if strings.HasSuffix(filename, ".pb.go") {
			pbGoFilesMap[dir] = append(pbGoFilesMap[dir], path)
		}

		return nil
	})
	if walkErr2 != nil {
		fmt.Println(walkErr2)
		os.Exit(1)
	}

	var err error
	for _, srcPbGoFiles := range pbGoFilesMap {
		for _, srcPbGoFile := range srcPbGoFiles {
			var dstPbGoFile string
			dstPbGoFile, err = filepath.Rel(modulePath, srcPbGoFile)
			if err != nil {
				fmt.Println(err)
				continue
			}
			err = os.Link(srcPbGoFile, dstPbGoFile)
			if err != nil {
				if os.IsNotExist(err) {
					fmt.Printf("'%s' does not exist\n", srcPbGoFile)
					continue
				}
				if os.IsPermission(err) {
					fmt.Println(err)
					continue
				}
				if os.IsExist(err) {
					err = os.Remove(dstPbGoFile)
					if err != nil {
						fmt.Printf("Failed to delete file '%s'\n", dstPbGoFile)
						continue
					}
					err = os.Rename(srcPbGoFile, dstPbGoFile)
					if err != nil {
						fmt.Printf("Can not move '%s' to '%s'\n", srcPbGoFile, dstPbGoFile)
					}
					continue
				}
			}
			err = os.Rename(srcPbGoFile, dstPbGoFile)
			if err != nil {
				fmt.Printf("Can not move '%s' to '%s'\n", srcPbGoFile, dstPbGoFile)
			}
			continue
		}
	}

	if err == nil {
		err = os.RemoveAll(strings.Split(modulePath, "/")[0])
		if err != nil {
			fmt.Println(err)
		}
	}
}
