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

func (this *TypedSettings) Load(message proto.Message) error {
	targetType := GetType(message)
	if targetType != this.Type {
		return errors.New("Have type " + this.Type + ", but retrieved for " + targetType)
	}
	return proto.Unmarshal(this.Settings, message)
}

func (this *TypedSettings) GetInstance() (interface{}, error) {
	instance, err := GetInstance(this.Type)
	if err != nil {
		return nil, err
	}
	if err := proto.Unmarshal(this.Settings, instance.(proto.Message)); err != nil {
		return nil, err
	}
	return instance, nil
}
