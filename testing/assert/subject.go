package assert

import (
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
