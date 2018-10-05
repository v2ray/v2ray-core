package confloader

import (
	"io"
	"os"
)

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
