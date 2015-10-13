package unit

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

func (subject *ErrorSubject) IsNotNil() {
	if subject.value == nil {
		subject.FailWithMessage("Not true that the error is not nil.")
	}
}
