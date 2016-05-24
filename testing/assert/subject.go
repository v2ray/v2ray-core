package assert

type Subject struct {
	disp string
	a    *Assert
}

func (subject *Subject) Fail(verb string, other string) {
	subject.FailWithMessage("Not true that " + subject.DisplayString() + " " + verb + " <" + other + ">.")
}

func (subject *Subject) FailWithMessage(message string) {
	subject.a.Fail(message)
}

func (subject *Subject) DisplayString() string {
	value := subject.disp
	if len(value) == 0 {
		value = "unknown"
	}
	return "<" + value + ">"
}
