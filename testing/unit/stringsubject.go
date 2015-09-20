package unit

type StringSubject struct {
	*Subject
	value string
}

func NewStringSubject(base *Subject, value string) *StringSubject {
	return &StringSubject{
		Subject: base,
		value:   value,
	}
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
