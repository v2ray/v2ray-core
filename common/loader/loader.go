package loader

import (
	"errors"
	"v2ray.com/core/common"
)

var (
	ErrUnknownConfigID = errors.New("Unknown config ID.")
)

type ConfigCreator func() interface{}

type ConfigCreatorCache map[string]ConfigCreator

func (this ConfigCreatorCache) RegisterCreator(id string, creator ConfigCreator) error {
	if _, found := this[id]; found {
		return common.ErrDuplicatedName
	}

	this[id] = creator
	return nil
}

func (this ConfigCreatorCache) CreateConfig(id string) (interface{}, error) {
	creator, found := this[id]
	if !found {
		return nil, ErrUnknownConfigID
	}
	return creator(), nil
}

type ConfigLoader interface {
	Load([]byte) (interface{}, string, error)
	LoadWithID([]byte, string) (interface{}, error)
}
