package serial

type StringTList []StringT

func NewStringTList(raw []string) *StringTList {
	list := StringTList(make([]StringT, len(raw)))
	for idx, str := range raw {
		list[idx] = StringT(str)
	}
	return &list
}

func (this *StringTList) Len() int {
	return len(*this)
}
