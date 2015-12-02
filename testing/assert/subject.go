package assert

import (
	"github.com/v2ray/v2ray-core/common/serial"
	v2testing "github.com/v2ray/v2ray-core/testing"
)

type Subject struct {
	name string
}

func NewSubject() *Subject {
	return &Subject{
		name: "",
	}
}

func (subject *Subject) Fail(displayString string, verb string, other serial.String) {
	subject.FailWithMessage("Not true that " + displayString + " " + verb + " <" + other.String() + ">.")
}

func (subject *Subject) FailWithMessage(message string) {
	v2testing.Fail(message)
}

func (subject *Subject) Named(name string) {
	subject.name = name
}

func (subject *Subject) DisplayString(value string) string {
	if len(value) == 0 {
		value = "unknown"
	}
	if len(subject.name) == 0 {
		return "<" + value + ">"
	}
	return subject.name + "(<" + value + ">)"
}
