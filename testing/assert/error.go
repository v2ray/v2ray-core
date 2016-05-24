package assert

func (this *Assert) Error(value error) *ErrorSubject {
	valueStr := ""
	if value != nil {
		valueStr = value.Error()
	}
	return &ErrorSubject{
		Subject: Subject{
			a:    this,
			disp: valueStr,
		},
		value: value,
	}
}

type ErrorSubject struct {
	Subject
	value error
}

func (subject *ErrorSubject) Equals(expectation error) {
	if subject.value != expectation {
		subject.Fail("is equal to", expectation.Error())
	}
}

func (subject *ErrorSubject) IsNil() {
	if subject.value != nil {
		subject.Fail("is", "nil")
	}
}

func (subject *ErrorSubject) IsNotNil() {
	if subject.value == nil {
		subject.Fail("is not", "nil")
	}
}
