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

// An implemention of Matcher
// Visible for testing only.
type FullGroupMatcher struct {
	matchers map[string]bool
}

func (g *FullGroupMatcher) New() {
	g.matchers = make(map[string]bool)
}

func (g *FullGroupMatcher) Add(domain string) {
	g.matchers[domain] = true
}

func (g *FullGroupMatcher) addMatcher(m fullMatcher) {
	g.Add(string(m))
}

func (g *FullGroupMatcher) Match(str string) bool {
	return g.matchers[str]
}
