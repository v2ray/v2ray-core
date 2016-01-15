package serial

type StringLiteralList []StringLiteral

func NewStringLiteralList(raw []string) *StringLiteralList {
	list := StringLiteralList(make([]StringLiteral, len(raw)))
	for idx, str := range raw {
		list[idx] = StringLiteral(str)
	}
	return &list
}

func (this *StringLiteralList) Len() int {
	return len(*this)
}
