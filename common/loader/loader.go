package loader

import (
	"errors"
)

var (
	ErrConfigIDKeyNotFound = errors.New("Config ID key is not found.")
	ErrConfigIDExists      = errors.New("Config ID already exists.")
	ErrUnknownConfigID     = errors.New("Unknown config ID.")
)

type ConfigCreator func() interface{}

type ConfigLoader interface {
	RegisterCreator(string, ConfigCreator) error
	CreateConfig(string) (interface{}, error)
	Load([]byte) (interface{}, error)
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
		return ErrConfigIDExists
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
