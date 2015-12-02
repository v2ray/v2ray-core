package serial

type String interface {
	String() string
}

type StringLiteral string

func (this StringLiteral) String() string {
	return string(this)
}
