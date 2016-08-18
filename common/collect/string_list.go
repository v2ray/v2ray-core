package collect

type StringList []string

func NewStringList(raw []string) *StringList {
	list := StringList(raw)
	return &list
}

func (this StringList) Len() int {
	return len(this)
}
