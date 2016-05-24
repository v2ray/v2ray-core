package assert

import (
	"strings"
)

func (this *Assert) String(value string) *StringSubject {
	return &StringSubject{
		Subject: Subject{
			a:    this,
			disp: value,
		},
		value: value,
	}
}

type StringSubject struct {
	Subject
	value string
}

func (subject *StringSubject) Equals(expectation string) {
	if subject.value != expectation {
		subject.Fail("is equal to", expectation)
	}
}

func (subject *StringSubject) NotEquals(expectation string) {
	if subject.value == expectation {
		subject.Fail("is not equal to ", expectation)
	}
}

func (subject *StringSubject) Contains(substring string) {
	if !strings.Contains(subject.value, substring) {
		subject.Fail("contains", substring)
	}
}

func (subject *StringSubject) NotContains(substring string) {
	if strings.Contains(subject.value, substring) {
		subject.Fail("doesn't contain", substring)
	}
}
