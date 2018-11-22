package external

import (
	"io"
	"os"
	"strings"

	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/platform/ctlcmd"
	"v2ray.com/core/main/confloader"
)

//go:generate errorgen

func loadConfigFile(configFile string) (io.ReadCloser, error) {
	if configFile == "stdin:" {
		return os.Stdin, nil
	}

	if strings.HasPrefix(configFile, "http://") || strings.HasPrefix(configFile, "https://") {
		content, err := ctlcmd.Run([]string{"fetch", configFile}, nil)
		if err != nil {
			return nil, err
		}
		return &buf.MultiBufferContainer{
			MultiBuffer: content,
		}, nil
	}

	fixedFile := os.ExpandEnv(configFile)
	file, err := os.Open(fixedFile)
	if err != nil {
		return nil, newError("config file not readable").Base(err)
	}
	defer file.Close()

	content, err := buf.ReadFrom(file)
	if err != nil {
		return nil, newError("failed to load config file: ", fixedFile).Base(err).AtWarning()
	}
	return &buf.MultiBufferContainer{
		MultiBuffer: content,
	}, nil
}

func init() {
	confloader.EffectiveConfigFileLoader = loadConfigFile
}
