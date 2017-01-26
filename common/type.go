package common

import (
	"context"
	"errors"
	"reflect"
)

type creator func(ctx context.Context, config interface{}) (interface{}, error)

var (
	typeCreatorRegistry = make(map[reflect.Type]creator)
)

func RegisterConfig(config interface{}, configCreator creator) error {
	configType := reflect.TypeOf(config)
	if _, found := typeCreatorRegistry[configType]; found {
		return errors.New("Common: " + configType.Name() + " is already registered.")
	}
	typeCreatorRegistry[configType] = configCreator
	return nil
}

func CreateObject(ctx context.Context, config interface{}) (interface{}, error) {
	configType := reflect.TypeOf(config)
	creator, found := typeCreatorRegistry[configType]
	if !found {
		return nil, errors.New("Common: " + configType.String() + " is not registered.")
	}
	return creator(ctx, config)
}
