package assert

import (
	"strings"

	"github.com/v2ray/v2ray-core/common/serial"
)

func StringLiteral(value string) *StringSubject {
	return String(serial.StringT((value)))
}

func String(value serial.String) *StringSubject {
	return &StringSubject{value: value}
}

type StringSubject struct {
	Subject
	value serial.String
}

func (subject *StringSubject) Named(name string) *StringSubject {
	subject.Subject.Named(name)
	return subject
}

func (subject *StringSubject) DisplayString() string {
	return subject.Subject.DisplayString(subject.value.String())
}

func (subject *StringSubject) Equals(expectation string) {
	if subject.value.String() != expectation {
		subject.Fail(subject.DisplayString(), "is equal to", serial.StringT(expectation))
	}
}

func (subject *StringSubject) NotEquals(expectation string) {
	if subject.value.String() == expectation {
		subject.Fail(subject.DisplayString(), "is not equal to ", serial.StringT(expectation))
	}
}

func (subject *StringSubject) Contains(substring serial.String) {
	if !strings.Contains(subject.value.String(), substring.String()) {
		subject.Fail(subject.DisplayString(), "contains", substring)
	}
}

func (subject *StringSubject) NotContains(substring serial.String) {
	if strings.Contains(subject.value.String(), substring.String()) {
		subject.Fail(subject.DisplayString(), "doesn't contain", substring)
	}
}
