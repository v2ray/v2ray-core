package loader

import (
	"errors"
	"reflect"

	"github.com/golang/protobuf/proto"
)

func NewTypedSettings(message proto.Message) *TypedSettings {
	if message == nil {
		return nil
	}
	settings, _ := proto.Marshal(message)
	return &TypedSettings{
		Type:     GetType(message),
		Settings: settings,
	}
}

func GetType(message proto.Message) string {
	return proto.MessageName(message)
}

func GetInstance(messageType string) (interface{}, error) {
	mType := proto.MessageType(messageType).Elem()
	if mType == nil {
		return nil, errors.New("Unknown type: " + messageType)
	}
	return reflect.New(mType).Interface(), nil
}

func (v *TypedSettings) Load(message proto.Message) error {
	targetType := GetType(message)
	if targetType != v.Type {
		return errors.New("Have type " + v.Type + ", but retrieved for " + targetType)
	}
	return proto.Unmarshal(v.Settings, message)
}

func (v *TypedSettings) GetInstance() (interface{}, error) {
	instance, err := GetInstance(v.Type)
	if err != nil {
		return nil, err
	}
	if err := proto.Unmarshal(v.Settings, instance.(proto.Message)); err != nil {
		return nil, err
	}
	return instance, nil
}
