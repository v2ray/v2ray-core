package serial

import (
	"errors"
	"reflect"

	"github.com/golang/protobuf/proto"
)

// ToTypedMessage converts a proto Message into TypedMessage.
func ToTypedMessage(message proto.Message) *TypedMessage {
	if message == nil {
		return nil
	}
	settings, _ := proto.Marshal(message)
	return &TypedMessage{
		Type:  GetMessageType(message),
		Value: settings,
	}
}

// GetMessageType returns the name of this proto Message.
func GetMessageType(message proto.Message) string {
	return proto.MessageName(message)
}

// GetInstance creates a new instance of the message with messageType.
func GetInstance(messageType string) (interface{}, error) {
	mType := proto.MessageType(messageType)
	if mType == nil || mType.Elem() == nil {
		return nil, errors.New("Serial: Unknown type: " + messageType)
	}
	return reflect.New(mType.Elem()).Interface(), nil
}

// GetInstance converts current TypedMessage into a proto Message.
func (v *TypedMessage) GetInstance() (proto.Message, error) {
	instance, err := GetInstance(v.Type)
	if err != nil {
		return nil, err
	}
	protoMessage := instance.(proto.Message)
	if err := proto.Unmarshal(v.Value, protoMessage); err != nil {
		return nil, err
	}
	return protoMessage, nil
}
