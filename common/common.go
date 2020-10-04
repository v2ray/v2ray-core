// Package common contains common utilities that are shared among other packages.
// See each sub-package for detail.
package common

import (
	"fmt"
	"go/build"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"v2ray.com/core/common/errors"
)

//go:generate go run v2ray.com/core/common/errors/errorgen

var (
	// ErrNoClue is for the situation that existing information is not enough to make a decision. For example, Router may return this error when there is no suitable route.
	ErrNoClue = errors.New("not enough information for making a decision")
)

// Must panics if err is not nil.
func Must(err error) {
	if err != nil {
		panic(err)
	}
}

// Must2 panics if the second parameter is not nil, otherwise returns the first parameter.
func Must2(v interface{}, err error) interface{} {
	Must(err)
	return v
}

// Error2 returns the err from the 2nd parameter.
func Error2(v interface{}, err error) error {
	return err
}

// envFile returns the name of the Go environment configuration file.
// Copy from https://github.com/golang/go/blob/c4f2a9788a7be04daf931ac54382fbe2cb754938/src/cmd/go/internal/cfg/cfg.go#L150-L166
func envFile() (string, error) {
	if file := os.Getenv("GOENV"); file != "" {
		if file == "off" {
			return "", fmt.Errorf("GOENV=off")
		}
		return file, nil
	}
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	if dir == "" {
		return "", fmt.Errorf("missing user-config dir")
	}
	return filepath.Join(dir, "go", "env"), nil
}

// GetRuntimeEnv returns the value of runtime environment variable,
// that is set by running following command: `go env -w key=value`.
func GetRuntimeEnv(key string) (string, error) {
	file, err := envFile()
	if err != nil {
		return "", err
	}
	if file == "" {
		return "", fmt.Errorf("missing runtime env file")
	}
	var data []byte
	var runtimeEnv string
	data, readErr := ioutil.ReadFile(file)
	if readErr != nil {
		return "", readErr
	}
	envStrings := strings.Split(string(data), "\n")
	for _, envItem := range envStrings {
		envItem = strings.TrimSuffix(envItem, "\r")
		envKeyValue := strings.Split(envItem, "=")
		if strings.EqualFold(strings.TrimSpace(envKeyValue[0]), key) {
			runtimeEnv = strings.TrimSpace(envKeyValue[1])
		}
	}
	return runtimeEnv, nil
}

// GetGOBIN returns GOBIN environment variable as a string. It will NOT be empty.
func GetGOBIN() string {
	// The one set by user explicitly by `export GOBIN=/path` or `env GOBIN=/path command`
	GOBIN := os.Getenv("GOBIN")
	if GOBIN == "" {
		var err error
		// The one set by user by running `go env -w GOBIN=/path`
		GOBIN, err = GetRuntimeEnv("GOBIN")
		if err != nil {
			// The default one that Golang uses
			return filepath.Join(build.Default.GOPATH, "bin")
		}
		if GOBIN == "" {
			return filepath.Join(build.Default.GOPATH, "bin")
		}
		return GOBIN
	}
	return GOBIN
}

// GetGOPATH returns GOPATH environment variable as a string. It will NOT be empty.
func GetGOPATH() string {
	// The one set by user explicitly by `export GOPATH=/path` or `env GOPATH=/path command`
	GOPATH := os.Getenv("GOPATH")
	if GOPATH == "" {
		var err error
		// The one set by user by running `go env -w GOPATH=/path`
		GOPATH, err = GetRuntimeEnv("GOPATH")
		if err != nil {
			// The default one that Golang uses
			return build.Default.GOPATH
		}
		if GOPATH == "" {
			return build.Default.GOPATH
		}
		return GOPATH
	}
	return GOPATH
}

// GetModuleName returns the value of module in `go.mod` file.
func GetModuleName(pathToProjectRoot string) (string, error) {
	var moduleName string
	loopPath := pathToProjectRoot
	for {
		if idx := strings.LastIndex(loopPath, string(filepath.Separator)); idx >= 0 {
			gomodPath := filepath.Join(loopPath, "go.mod")
			gomodBytes, err := ioutil.ReadFile(gomodPath)
			if err != nil {
				loopPath = loopPath[:idx]
				continue
			}

			gomodContent := string(gomodBytes)
			moduleIdx := strings.Index(gomodContent, "module ")
			newLineIdx := strings.Index(gomodContent, "\n")

			if moduleIdx >= 0 {
				if newLineIdx >= 0 {
					moduleName = strings.TrimSpace(gomodContent[moduleIdx+6 : newLineIdx])
					moduleName = strings.TrimSuffix(moduleName, "\r")
				} else {
					moduleName = strings.TrimSpace(gomodContent[moduleIdx+6:])
				}
				return moduleName, nil
			}
			return "", fmt.Errorf("can not get module path in `%s`", gomodPath)
		}
		break
	}
	return moduleName, fmt.Errorf("no `go.mod` file in every parent directory of `%s`", pathToProjectRoot)
}
