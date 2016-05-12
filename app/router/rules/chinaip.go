package rules

//go:generate go run chinaip_gen.go

func NewChinaIPRule(tag string) *Rule {
	return &Rule{
		Tag:       tag,
		Condition: NewIPv4Matcher(chinaIPNet),
	}
}
