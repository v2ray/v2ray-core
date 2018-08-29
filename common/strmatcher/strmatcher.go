package strmatcher

import (
	"regexp"
	"sync"
	"time"

	"v2ray.com/core/common"
	"v2ray.com/core/common/task"
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

type cacheEntry struct {
	timestamp time.Time
	result    uint32
}

// CachedMatcherGroup is a IndexMatcher with cachable results.
type CachedMatcherGroup struct {
	sync.RWMutex
	group   *MatcherGroup
	cache   map[string]cacheEntry
	cleanup *task.Periodic
}

// NewCachedMatcherGroup creats a new CachedMatcherGroup.
func NewCachedMatcherGroup(g *MatcherGroup) *CachedMatcherGroup {
	r := &CachedMatcherGroup{
		group: g,
		cache: make(map[string]cacheEntry),
	}
	r.cleanup = &task.Periodic{
		Interval: time.Second * 30,
		Execute: func() error {
			r.Lock()
			defer r.Unlock()

			if len(r.cache) == 0 {
				return nil
			}

			expire := time.Now().Add(-1 * time.Second * 120)
			for p, e := range r.cache {
				if e.timestamp.Before(expire) {
					delete(r.cache, p)
				}
			}

			if len(r.cache) == 0 {
				r.cache = make(map[string]cacheEntry)
			}

			return nil
		},
	}
	common.Must(r.cleanup.Start())
	return r
}

// Match implements IndexMatcher.Match.
func (g *CachedMatcherGroup) Match(pattern string) uint32 {
	g.RLock()
	r, f := g.cache[pattern]
	g.RUnlock()
	if f {
		return r.result
	}

	mr := g.group.Match(pattern)

	g.Lock()
	g.cache[pattern] = cacheEntry{
		result:    mr,
		timestamp: time.Now(),
	}
	g.Unlock()

	return mr
}
