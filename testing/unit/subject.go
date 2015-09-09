package unit

type Subject struct {
	assert *Assertion
	name   string
}

func NewSubject(assert *Assertion) *Subject {
	subject := new(Subject)
	subject.assert = assert
	subject.name = ""
	return subject
}

func (subject *Subject) FailWithMessage(message string) {
	subject.assert.t.Error(message)
}

func (subject *Subject) Named(name string) {
	subject.name = name
}

func (subject *Subject) DisplayString(value string) string {
	if len(value) == 0 {
		value = "unknown"
	}
	if len(subject.name) == 0 {
		return "<" + value + ">"
	}
	return subject.name + "(<" + value + ">)"
}
