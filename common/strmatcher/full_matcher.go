package strmatcher

// FullMatcherGroup is an implemention of IndexMatcher
type FullMatcherGroup struct {
	matchers map[string]uint32
}

// Add a domain for matching
func (g *FullMatcherGroup) Add(domain string, value uint32) {
	if g.matchers == nil {
		g.matchers = make(map[string]uint32)
	}

	g.matchers[domain] = value
}

func (g *FullMatcherGroup) addMatcher(m fullMatcher, value uint32) {
	g.Add(string(m), value)
}

// Match is an implementation of IndexMatcher.Match.
func (g *FullMatcherGroup) Match(str string) uint32 {
	if g.matchers == nil {
		return 0
	}

	return g.matchers[str]
}

// FullGroupMatcher is an implemention of Matcher
// Visible for testing only.
type FullGroupMatcher struct {
	matchers map[string]bool
}

// New a FullGroupMatcher
func (g *FullGroupMatcher) New() {
	g.matchers = make(map[string]bool)
}

// Add a domain for matching
func (g *FullGroupMatcher) Add(domain string) {
	g.matchers[domain] = true
}

func (g *FullGroupMatcher) addMatcher(m fullMatcher) {
	g.Add(string(m))
}

// Match is an implementation of Matcher.Match.
func (g *FullGroupMatcher) Match(str string) bool {
	return g.matchers[str]
}
