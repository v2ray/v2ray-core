package router

import (
	"errors"
)

type ConfigObjectCreator func([]byte) (interface{}, error)

var (
	configCache map[string]ConfigObjectCreator

	ErrorRouterNotFound = errors.New("Router not found.")
)

func RegisterRouterConfig(strategy string, creator ConfigObjectCreator) error {
	// TODO: check strategy
	configCache[strategy] = creator
	return nil
}

func CreateRouterConfig(strategy string, data []byte) (interface{}, error) {
	creator, found := configCache[strategy]
	if !found {
		return nil, ErrorRouterNotFound
	}
	return creator(data)
}

func init() {
	configCache = make(map[string]ConfigObjectCreator)
}
