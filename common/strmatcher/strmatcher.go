package strmatcher

import "regexp"

type Matcher interface {
	Match(string) bool
}

type Type byte

const (
	Full Type = iota
	Substr
	Domain
	Regex
)

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

type matcherEntry struct {
	m  Matcher
	id uint32
}

type MatcherGroup struct {
	count         uint32
	fullMatchers  map[string]uint32
	otherMatchers []matcherEntry
}

func NewMatcherGroup() *MatcherGroup {
	return &MatcherGroup{
		count:        1,
		fullMatchers: make(map[string]uint32),
	}
}

func (g *MatcherGroup) Add(m Matcher) uint32 {
	c := g.count
	g.count++

	if fm, ok := m.(fullMatcher); ok {
		g.fullMatchers[string(fm)] = c
	} else {
		g.otherMatchers = append(g.otherMatchers, matcherEntry{
			m:  m,
			id: c,
		})
	}

	return c
}

func (g *MatcherGroup) Match(pattern string) uint32 {
	if c, f := g.fullMatchers[pattern]; f {
		return c
	}

	for _, e := range g.otherMatchers {
		if e.m.Match(pattern) {
			return e.id
		}
	}

	return 0
}

func (g *MatcherGroup) Size() uint32 {
	return g.count
}
