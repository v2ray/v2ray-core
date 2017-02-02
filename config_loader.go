package core

import (
	"io"
	"io/ioutil"

	"github.com/golang/protobuf/proto"

	"v2ray.com/core/common/errors"
)

type ConfigLoader func(input io.Reader) (*Config, error)

var configLoaderCache = make(map[ConfigFormat]ConfigLoader)

func RegisterConfigLoader(format ConfigFormat, loader ConfigLoader) error {
	configLoaderCache[format] = loader
	return nil
}

func LoadConfig(format ConfigFormat, input io.Reader) (*Config, error) {
	loader, found := configLoaderCache[format]
	if !found {
		return nil, errors.New("Core: ", ConfigFormat_name[int32(format)], " is not loadable.")
	}
	return loader(input)
}

func LoadProtobufConfig(input io.Reader) (*Config, error) {
	config := new(Config)
	data, _ := ioutil.ReadAll(input)
	if err := proto.Unmarshal(data, config); err != nil {
		return nil, err
	}
	return config, nil
}

func init() {
	RegisterConfigLoader(ConfigFormat_Protobuf, LoadProtobufConfig)
}
