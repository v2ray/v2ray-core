package loader

import (
	"errors"
	"v2ray.com/core/common"
)

var (
	ErrUnknownConfigID = errors.New("Unknown config ID.")
)

type ConfigCreator func() interface{}

type ConfigLoader interface {
	RegisterCreator(string, ConfigCreator) error
	CreateConfig(string) (interface{}, error)
	Load([]byte) (interface{}, string, error)
	LoadWithID([]byte, string) (interface{}, error)
}

type BaseConfigLoader struct {
	creators map[string]ConfigCreator
}

func NewBaseConfigLoader() *BaseConfigLoader {
	return &BaseConfigLoader{
		creators: make(map[string]ConfigCreator),
	}
}

func (this *BaseConfigLoader) RegisterCreator(id string, creator ConfigCreator) error {
	if _, found := this.creators[id]; found {
		return common.ErrDuplicatedName
	}

	this.creators[id] = creator
	return nil
}

func (this *BaseConfigLoader) CreateConfig(id string) (interface{}, error) {
	creator, found := this.creators[id]
	if !found {
		return nil, ErrUnknownConfigID
	}
	return creator(), nil
}
