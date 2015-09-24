package unit

import (
	"fmt"

	"github.com/v2ray/v2ray-core/common/errors"
)

type ErrorSubject struct {
	*Subject
	value error
}

func NewErrorSubject(base *Subject, value error) *ErrorSubject {
	return &ErrorSubject{
		Subject: base,
		value:   value,
	}
}

func (subject *ErrorSubject) Named(name string) *ErrorSubject {
	subject.Subject.Named(name)
	return subject
}

func (subject *ErrorSubject) Fail(verb string, other error) {
	subject.FailWithMessage("Not true that " + subject.DisplayString() + " " + verb + " <" + other.Error() + ">.")
}

func (subject *ErrorSubject) DisplayString() string {
	return subject.Subject.DisplayString(subject.value.Error())
}

func (subject *ErrorSubject) Equals(expectation error) {
	if subject.value != expectation {
		subject.Fail("is equal to", expectation)
	}
}

func (subject *ErrorSubject) IsNil() {
	if subject.value != nil {
		subject.FailWithMessage("Not true that " + subject.DisplayString() + " is nil.")
	}
}

func (subject *ErrorSubject) HasCode(code int) {
	if !errors.HasCode(subject.value, code) {
		subject.FailWithMessage(fmt.Sprintf("Not ture that %s has error code 0x%04X.", subject.DisplayString(), code))
	}
}
