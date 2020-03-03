package strmatcher

import (
	"regexp"
	"strings"
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

// NewSubstrMatcher creates a new substr matcher
// For app/stats/command/command.go only
func NewSubstrMatcher(pattern string) Matcher {
	return substrMatcher(pattern)
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
func (mg *MatcherGroup) Add(m Matcher) {
	c := mg.count
	switch tm := m.(type) {
	case fullMatcher:
		mg.fullMatcher.addMatcher(tm, c)
	case domainMatcher:
		mg.domainMatcher.addMatcher(tm, c)
	default:
		mg.otherMatchers = append(mg.otherMatchers, matcherEntry{
			m:  m,
			id: c,
		})
	}
}

type groupMatcher interface {
	Add(m Matcher)
}

type patternParser struct {
	unparsedNumber int
	nextPattern    int
	external       map[string][]string
}

// Parse a pattern to a part of MatcherGroup
func (p *patternParser) subPattern(mg groupMatcher, pattern string) error {
	cmd := pattern[0]
	length := len(pattern)
	p.nextPattern = length
	// For the matchers which have not child matcher
	if p.unparsedNumber != 0 && cmd != '&' && cmd != '!' {
		pos := strings.IndexByte(pattern, ' ')
		if pos != -1 {
			p.nextPattern = pos
		}
	}
	left := pattern[1:p.nextPattern]
	var m Matcher
	switch cmd {
	case 'd':
		// Domain
		m = domainMatcher(left)
	case 'r':
		// Regexp
		// Return error at the end of function
		r, err := regexp.Compile(left)
		if err != nil {
			return err
		}
		m = &regexMatcher{
			pattern: r,
		}
	case 'k':
		// Keyword
		m = substrMatcher(left)
	case 'f':
		// Full
		m = fullMatcher(left)
	case 'e':
		// External
		sp := &patternParser{
			unparsedNumber: 0,
			nextPattern:    0,
			external:       p.external,
		}
		// Use extra parser to avoid unnecessary calculation
		for _, newPattern := range p.external[left] {
			sp.subPattern(mg, newPattern)
		}
	case '!':
		// Not
		smg := NewOrMatcher()
		p.subPattern(smg, pattern[1:])
		p.nextPattern++
		m = &notMatcher{
			matcher: smg,
		}
	case '&':
		a := NewOrMatcher()
		b := NewOrMatcher()
		p.unparsedNumber++
		err := p.subPattern(a, pattern[1:])
		if err != nil {
			return err
		}
		lenA := p.nextPattern
		err = p.subPattern(b, pattern[lenA+2:])
		if err != nil {
			return err
		}
		p.unparsedNumber--
		p.nextPattern += lenA + 2
		m = &andMatcher{
			matcherA: a,
			matcherB: b,
		}
	case '|':
		p.unparsedNumber++
		err := p.subPattern(mg, pattern[1:])
		if err != nil {
			return err
		}
		lenA := p.nextPattern
		err = p.subPattern(mg, pattern[lenA+2:])
		if err != nil {
			return err
		}
		p.unparsedNumber--
		p.nextPattern += lenA + 2
	default:
		panic("Unknown type")
	}
	if m != nil {
		mg.Add(m)
	}
	return nil
}

// ParsePattern parses a pattern to a part of MatcherGroup and return its index. The index will never be 0.
func (mg *MatcherGroup) ParsePattern(pattern string, extern map[string][]string) (uint32, error) {
	mg.count++
	p := &patternParser{
		unparsedNumber: 0,
		nextPattern:    0,
		external:       extern,
	}
	return mg.count, p.subPattern(mg, pattern)
}

// Match implements IndexMatcher.Match.
func (mg *MatcherGroup) Match(pattern string) uint32 {
	if c := mg.fullMatcher.Match(pattern); c > 0 {
		return c
	}

	if c := mg.domainMatcher.Match(pattern); c > 0 {
		return c
	}

	for _, e := range mg.otherMatchers {
		if e.m.Match(pattern) {
			return e.id
		}
	}

	return 0
}

// OrMatcher is a implementation of Matcher
type OrMatcher struct {
	fullMatchers   FullGroupMatcher
	domainMatchers DomainGroupMatcher
	otherMatchers  []Matcher
}

// NewOrMatcher creates an OrMatcher
func NewOrMatcher() (g *OrMatcher) {
	g = new(OrMatcher)
	g.fullMatchers.New()
	return
}

// Match implements Matcher.Match.
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

// ParsePattern parses a pattern to a part of OrMatcher
func (g *OrMatcher) ParsePattern(pattern string, extern map[string][]string) error {
	p := &patternParser{
		unparsedNumber: 0,
		nextPattern:    0,
		external:       extern,
	}
	return p.subPattern(g, pattern)
}
