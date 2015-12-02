package assert

import (
	"strconv"
)

func Int(value int) *IntSubject {
	return &IntSubject{value: value}
}

type IntSubject struct {
	*Subject
	value int
}

func (subject *IntSubject) Named(name string) *IntSubject {
	subject.Subject.Named(name)
	return subject
}

func (subject *IntSubject) Fail(verb string, other int) {
	subject.FailWithMessage("Not true that " + subject.DisplayString() + " " + verb + " <" + strconv.Itoa(other) + ">.")
}

func (subject *IntSubject) DisplayString() string {
	return subject.Subject.DisplayString(strconv.Itoa(subject.value))
}

func (subject *IntSubject) Equals(expectation int) {
	if subject.value != expectation {
		subject.Fail("is equal to", expectation)
	}
}

func (subject *IntSubject) GreaterThan(expectation int) {
	if subject.value <= expectation {
		subject.Fail("is greater than", expectation)
	}
}

func (subject *IntSubject) LessThan(expectation int) {
	if subject.value >= expectation {
		subject.Fail("is less than", expectation)
	}
}
