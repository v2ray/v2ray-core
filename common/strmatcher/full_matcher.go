package strmatcher

// A implemention of IndexMatcher
type FullMatcherGroup struct {
	matchers map[string]uint32
}

func (g *FullMatcherGroup) Add(domain string, value uint32) {
	if g.matchers == nil {
		g.matchers = make(map[string]uint32)
	}

	g.matchers[domain] = value
}

func (g *FullMatcherGroup) addMatcher(m fullMatcher, value uint32) {
	g.Add(string(m), value)
}

func (g *FullMatcherGroup) Match(str string) uint32 {
	if g.matchers == nil {
		return 0
	}

	return g.matchers[str]
}

var fullExists = struct{}{}

// A implemention of Matcher
type fullGroupMatcher struct {
	matchers map[string]struct{}
}

func (g *fullGroupMatcher) New() {
	g.matchers = make(map[string]struct{})
}

func (g *fullGroupMatcher) Add(domain string) {
	g.matchers[domain] = fullExists
}

func (g *fullGroupMatcher) addMatcher(m fullMatcher) {
	g.Add(string(m))
}

func (g *fullGroupMatcher) Match(str string) bool {
	if len(g.matchers) == 0 {
		return false
	}
	_, exist := g.matchers[str]
	return exist
}
