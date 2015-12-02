package assert

import (
	"fmt"
)

func Pointer(value interface{}) *PointerSubject {
	return &PointerSubject{value: value}
}

type PointerSubject struct {
	Subject
	value interface{}
}

func (subject *PointerSubject) Named(name string) *PointerSubject {
	subject.Subject.Named(name)
	return subject
}

func (subject *PointerSubject) Fail(verb string, other interface{}) {
	subject.FailWithMessage(fmt.Sprintf("Not true that %s %s <%v>.", subject.DisplayString(), verb, other))
}

func (subject *PointerSubject) DisplayString() string {
	return subject.Subject.DisplayString(fmt.Sprintf("%v", subject.value))
}

func (subject *PointerSubject) Equals(expectation interface{}) {
	if subject.value != expectation {
		subject.Fail("is equal to", expectation)
	}
}

func (subject *PointerSubject) IsNil() {
	if subject.value != nil {
		subject.FailWithMessage("Not true that " + subject.DisplayString() + " is nil.")
	}
}

func (subject *PointerSubject) IsNotNil() {
	if subject.value == nil {
		subject.FailWithMessage("Not true that " + subject.DisplayString() + " is not nil.")
	}
}
