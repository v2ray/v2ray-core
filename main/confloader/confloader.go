package confloader

import (
	"io"
	"os"
)

//go:generate go run $GOPATH/src/v2ray.com/core/common/errors/errorgen/main.go -pkg confloader -path Main,ConfLoader

type configFileLoader func(string) (io.ReadCloser, error)

var (
	EffectiveConfigFileLoader configFileLoader
)

func LoadConfig(file string) (io.ReadCloser, error) {
	if EffectiveConfigFileLoader == nil {
		return os.Stdin, nil
	}
	return EffectiveConfigFileLoader(file)
}
