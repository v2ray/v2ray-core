package serial

// An interface for any objects that has string presentation.
type String interface {
	String() string
}

type StringLiteral string

func (this StringLiteral) String() string {
	return string(this)
}
