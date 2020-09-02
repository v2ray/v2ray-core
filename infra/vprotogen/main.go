package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"v2ray.com/core/common"
)

var (
	reporoot = flag.String("repo", "", "Repo for protobuf generation, such as v2ray.com/core")
)

func main() {
	flag.Parse()

	protofiles := make(map[string][]string)
	// protoc := protocMap[runtime.GOOS]
	protoc := "protoc"

	tmpRootDir, err := ioutil.TempDir("", "vprotogen")
	defer os.RemoveAll(tmpRootDir)
	common.Must(err)

	tmpDir := filepath.Join(tmpRootDir, "v2ray.com", "core")
	os.Mkdir(tmpDir, os.ModeDir)

	common.Must(CopyDir(*reporoot, tmpDir, func(filename string) bool {
		return strings.HasSuffix(filename, ".proto")
	}))

	filepath.Walk(tmpDir, func(path string, info os.FileInfo, err error) error {
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
		args := []string{"--proto_path", tmpRootDir, "--go_out", "plugins=grpc:" + tmpRootDir}
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

	common.Must(filepath.Walk(tmpRootDir, func(path string, info os.FileInfo, err error) error {
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
		content = bytes.Replace(content, []byte("\"golang.org/x/net/context\""), []byte("\"context\""), 1)
		// content = bytes.ReplaceAll(content, []byte("github_com_BattleRoach_v2ray_core_"), []byte("v2ray_com_core_"))

		pos := bytes.Index(content, []byte("\npackage"))
		if pos > 0 {
			content = content[pos+1:]
		}

		return ioutil.WriteFile(path, content, info.Mode())
	}))
	common.Must(CopyDir(tmpDir, *reporoot, func(filename string) bool {
		return strings.HasSuffix(filename, ".pb.go")
	}))
}

// File copies a single file from src to dst
func CopyFile(src, dst string) error {
	var err error
	var srcfd *os.File
	var dstfd *os.File
	var srcinfo os.FileInfo

	if srcfd, err = os.Open(src); err != nil {
		return err
	}
	defer srcfd.Close()

	if dstfd, err = os.Create(dst); err != nil {
		return err
	}
	defer dstfd.Close()

	if _, err = io.Copy(dstfd, srcfd); err != nil {
		return err
	}
	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}
	return os.Chmod(dst, srcinfo.Mode())
}

// Dir copies a whole directory recursively
func CopyDir(src string, dst string, predicate func(filepath string) bool) error {
	var err error
	var fds []os.FileInfo
	var srcinfo os.FileInfo

	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}

	if err = os.MkdirAll(dst, srcinfo.Mode()); err != nil {
		return err
	}

	if fds, err = ioutil.ReadDir(src); err != nil {
		return err
	}
	for _, fd := range fds {
		srcfp := path.Join(src, fd.Name())
		dstfp := path.Join(dst, fd.Name())

		if fd.IsDir() {
			if err = CopyDir(srcfp, dstfp, predicate); err != nil {
				fmt.Println(err)
			}
		} else {
			if predicate(fd.Name()) {
				if err = CopyFile(srcfp, dstfp); err != nil {
					fmt.Println(err)
				}
			}
		}
	}
	return nil
}
