package strmatcher

import (
	"regexp"
)

// Matcher is the interface to determine a string matches a pattern.
type Matcher interface {
	// Match returns true if the given string matches a predefined pattern.
	Match(string) bool
}

// Type is the type of the matcher.
type Type byte

const (
	// Full is the type of matcher that the input string must exactly equal to the pattern.
	Full Type = iota
	// Substr is the type of matcher that the input string must contain the pattern as a sub-string.
	Substr
	// Domain is the type of matcher that the input string must be a sub-domain or itself of the pattern.
	Domain
	// Regex is the type of matcher that the input string must matches the regular-expression pattern.
	Regex
)

// New creates a new Matcher based on the given pattern.
func (t Type) New(pattern string) (Matcher, error) {
	switch t {
	case Full:
		return fullMatcher(pattern), nil
	case Substr:
		return substrMatcher(pattern), nil
	case Domain:
		return domainMatcher(pattern), nil
	case Regex:
		r, err := regexp.Compile(pattern)
		if err != nil {
			return nil, err
		}
		return &regexMatcher{
			pattern: r,
		}, nil
	default:
		panic("Unknown type")
	}
}

// IndexMatcher is the interface for matching with a group of matchers.
type IndexMatcher interface {
	// Match returns the the index of a matcher that matches the input. It returns 0 if no such matcher exists.
	Match(input string) uint32
}

type matcherEntry struct {
	m  Matcher
	id uint32
}

// MatcherGroup is an implementation of IndexMatcher.
// Empty initialization works.
type MatcherGroup struct {
	count         uint32
	fullMatcher   FullMatcherGroup
	domainMatcher DomainMatcherGroup
	otherMatchers []matcherEntry
}

// Add adds a new Matcher into the MatcherGroup without adding index
func (g *MatcherGroup) addChild(m Matcher) {
	c := g.count
	switch tm := m.(type) {
	case fullMatcher:
		g.fullMatcher.addMatcher(tm, c)
	case domainMatcher:
		g.domainMatcher.addMatcher(tm, c)
	default:
		g.otherMatchers = append(g.otherMatchers, matcherEntry{
			m:  m,
			id: c,
		})
	}
}

func (mg *MatcherGroup) subPattern(pattern string, extern map[string][]string) error {
	cmd := pattern[0]
	left := pattern[1:len(pattern)]
	var m Matcher = nil
	var err error = nil
	switch cmd {
	case 'd': // Domain
		m = domainMatcher(left)
	case 'r': // Regexp
		r, err := regexp.Compile(left)
		if err != nil {
			return nil
		}
		m = &regexMatcher{
			pattern: r,
		}
	case 'k': // Keyword
		m = substrMatcher(left)
	case 'f': // Full
		m = fullMatcher(left)
	case 'e': // External
		for _, newPattern := range extern[left] {
			mg.subPattern(newPattern, extern)
		}
	default:
		panic("Unknown type")
	}
	if m != nil {
		mg.addChild(m)
	}
	return err
}

func (mg *MatcherGroup) ParsePattern(pattern string, extern map[string][]string) (uint32, error) {
	mg.count++
	return mg.count, mg.subPattern(pattern, extern)
}

// Add adds a new Matcher into the MatcherGroup, and returns its index. The index will never be 0.
func (g *MatcherGroup) Add(m Matcher) uint32 {
	g.count++
	c := g.count

	switch tm := m.(type) {
	case fullMatcher:
		g.fullMatcher.addMatcher(tm, c)
	case domainMatcher:
		g.domainMatcher.addMatcher(tm, c)
	default:
		g.otherMatchers = append(g.otherMatchers, matcherEntry{
			m:  m,
			id: c,
		})
	}

	return c
}

// Match implements IndexMatcher.Match.
func (g *MatcherGroup) Match(pattern string) uint32 {
	if c := g.fullMatcher.Match(pattern); c > 0 {
		return c
	}

	if c := g.domainMatcher.Match(pattern); c > 0 {
		return c
	}

	for _, e := range g.otherMatchers {
		if e.m.Match(pattern) {
			return e.id
		}
	}

	return 0
}

// Size returns the number of matchers in the MatcherGroup.
func (g *MatcherGroup) Size() uint32 {
	return g.count
}

type OrMatcher struct {
	fullMatchers   fullGroupMatcher
	domainMatchers domainGroupMatcher
	otherMatchers  []Matcher
}

func (g *OrMatcher) New() {
	g.fullMatchers.New()
}

func (g *OrMatcher) Match(pattern string) bool {
	if g.fullMatchers.Match(pattern) || g.domainMatchers.Match(pattern) {
		return true
	}

	for _, e := range g.otherMatchers {
		if e.Match(pattern) {
			return true
		}
	}

	return false
}

// Add adds a new Matcher into the OrMatcher
func (g *OrMatcher) Add(m Matcher) {
	switch tm := m.(type) {
	case fullMatcher:
		g.fullMatchers.addMatcher(tm)
	case domainMatcher:
		g.domainMatchers.addMatcher(tm)
	default:
		g.otherMatchers = append(g.otherMatchers, m)
	}
}

func (g *OrMatcher) ParsePattern(pattern string, extern map[string][]string) error {
	cmd := pattern[0]
	left := pattern[1:len(pattern)]
	var m Matcher = nil
	var err error = nil
	switch cmd {
	case 'd': // Domain
		m = domainMatcher(left)
	case 'r': // Regexp
		r, err := regexp.Compile(left)
		if err != nil {
			return nil
		}
		m = &regexMatcher{
			pattern: r,
		}
	case 'k': // Keyword
		m = substrMatcher(left)
	case 'f': // Full
		m = fullMatcher(left)
	case 'e': // External
		for _, newPattern := range extern[left] {
			g.ParsePattern(newPattern, extern)
		}
	default:
		panic("Unknown type")
	}
	if m != nil {
		g.Add(m)
	}
	return err
}
