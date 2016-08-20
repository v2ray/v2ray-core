package assert

import (
	"reflect"

	"v2ray.com/core/common/serial"
)

func (this *Assert) Pointer(value interface{}) *PointerSubject {
	return &PointerSubject{
		Subject: Subject{
			a:    this,
			disp: serial.PointerToString(value),
		},
		value: value,
	}
}

type PointerSubject struct {
	Subject
	value interface{}
}

func (subject *PointerSubject) Equals(expectation interface{}) {
	if subject.value != expectation {
		subject.Fail("is equal to", serial.PointerToString(expectation))
	}
}

func (subject *PointerSubject) IsNil() {
	if subject.value == nil {
		return
	}

	valueType := reflect.TypeOf(subject.value)
	nilType := reflect.Zero(valueType)
	realValue := reflect.ValueOf(subject.value)

	if nilType != realValue {
		subject.Fail("is", "nil")
	}
}

func (subject *PointerSubject) IsNotNil() {
	if subject.value == nil {
		subject.Fail("is not", "nil")
	}

	valueType := reflect.TypeOf(subject.value)
	nilType := reflect.Zero(valueType)
	realValue := reflect.ValueOf(subject.value)

	if nilType == realValue {
		subject.Fail("is not", "nil")
	}
}
