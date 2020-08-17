package strmatcher

type FullMatcherGroup struct {
	matchers map[string][]uint32
}

func (g *FullMatcherGroup) Add(domain string, value uint32) {
	if g.matchers == nil {
		g.matchers = make(map[string][]uint32)
	}

	g.matchers[domain] = append(g.matchers[domain], value)
}

func (g *FullMatcherGroup) addMatcher(m fullMatcher, value uint32) {
	g.Add(string(m), value)
}

func (g *FullMatcherGroup) Match(str string) []uint32 {
	if g.matchers == nil {
		return nil
	}

	return g.matchers[str]
}
