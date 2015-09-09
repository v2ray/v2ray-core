package unit

import (
  "fmt"
)

type ErrorSubject struct {
	*Subject
	value error
}

func NewErrorSubject(base *Subject, value error) *ErrorSubject {
	subject := new(StringSubject)
	subject.Subject = base
	subject.value = value
	return subject
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
    subject.FailWithMethod("Not true that " + subject.DisplayString() + " is nil.")
  }
}