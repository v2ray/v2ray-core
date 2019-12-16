package confloader

import (
	"io"
	"os"
)

type configFileLoader func(string) (io.Reader, error)
type extconfigLoader func([]string) (io.Reader, error)

var (
	EffectiveConfigFileLoader configFileLoader
	EffectiveExtConfigLoader  extconfigLoader
)

func LoadConfig(file string) (io.Reader, error) {
	if EffectiveConfigFileLoader == nil {
		newError("external config module not loaded, reading from stdin").AtInfo().WriteToLog()
		return os.Stdin, nil
	}
	return EffectiveConfigFileLoader(file)
}

func LoadExtConfig(files []string) (io.Reader, error) {
	if EffectiveExtConfigLoader == nil {
		return nil, newError("external config module not loaded").AtError()
	}

	return EffectiveExtConfigLoader(files)
}
