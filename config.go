// +build !confonly

package core

import (
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/golang/protobuf/proto"
	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/cmdarg"
)

// ConfigFormat is a configurable format of V2Ray config file.
type ConfigFormat struct {
	Name      string
	Extension []string
	Loader    ConfigLoader
}

// ConfigLoader is a utility to load V2Ray config from external source.
type ConfigLoader func(input interface{}) (*Config, error)

var (
	configLoaderByName = make(map[string]*ConfigFormat)
	configLoaderByExt  = make(map[string]*ConfigFormat)
)

// RegisterConfigLoader add a new ConfigLoader.
func RegisterConfigLoader(format *ConfigFormat) error {
	name := strings.ToLower(format.Name)
	if _, found := configLoaderByName[name]; found {
		return newError(format.Name, " already registered.")
	}
	configLoaderByName[name] = format

	for _, ext := range format.Extension {
		lext := strings.ToLower(ext)
		if f, found := configLoaderByExt[lext]; found {
			return newError(ext, " already registered to ", f.Name)
		}
		configLoaderByExt[lext] = format
	}

	return nil
}

func getExtension(filename string) string {
	idx := strings.LastIndexByte(filename, '.')
	if idx == -1 {
		return ""
	}
	return filename[idx+1:]
}

// LoadConfig loads config with given format from given source.
// input accepts 2 different types:
// * []string slice of multiple filename/url(s) to open to read
// * io.Reader that reads a config content (the original way)
func LoadConfig(formatName string, filename string, input interface{}) (*Config, error) {
	ext := getExtension(filename)
	if len(ext) > 0 {
		if f, found := configLoaderByExt[ext]; found {
			return f.Loader(input)
		}
	}

	if f, found := configLoaderByName[formatName]; found {
		return f.Loader(input)
	}

	return nil, newError("Unable to load config in ", formatName).AtWarning()
}

func loadProtobufConfig(data []byte) (*Config, error) {
	config := new(Config)
	if err := proto.Unmarshal(data, config); err != nil {
		return nil, err
	}
	return config, nil
}

func init() {
	common.Must(RegisterConfigLoader(&ConfigFormat{
		Name:      "Protobuf",
		Extension: []string{"pb"},
		Loader: func(input interface{}) (*Config, error) {
			switch v := input.(type) {
			case cmdarg.Arg:
				if len(v) == 0 {
					return nil, newError("input has no element")
				}
				var data []byte
				var rerr error
				// pb type can only handle the first config
				if v[0] == "stdin:" {
					data, rerr = buf.ReadAllToBytes(os.Stdin)
				} else {
					data, rerr = ioutil.ReadFile(v[0])
				}
				common.Must(rerr)
				return loadProtobufConfig(data)
			case io.Reader:
				data, err := buf.ReadAllToBytes(v)
				common.Must(err)
				return loadProtobufConfig(data)
			default:
				return nil, newError("unknow type")
			}
		},
	}))
}
