package assert

func String(value string) *StringSubject {
	return &StringSubject{value: value}
}

type StringSubject struct {
	Subject
	value string
}

func (subject *StringSubject) Named(name string) *StringSubject {
	subject.Subject.Named(name)
	return subject
}

func (subject *StringSubject) Fail(verb string, other string) {
	subject.FailWithMessage("Not true that " + subject.DisplayString() + " " + verb + " <" + other + ">.")
}

func (subject *StringSubject) DisplayString() string {
	return subject.Subject.DisplayString(subject.value)
}

func (subject *StringSubject) Equals(expectation string) {
	if subject.value != expectation {
		subject.Fail("is equal to", expectation)
	}
}
