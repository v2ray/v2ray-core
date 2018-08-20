package strmatcher

import (
	"regexp"
	"sync"
	"time"

	"v2ray.com/core/common/task"
)

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

type IndexMatcher interface {
	Match(pattern string) uint32
}

type matcherEntry struct {
	m  Matcher
	id uint32
}

type MatcherGroup struct {
	count         uint32
	fullMatcher   FullMatcherGroup
	domainMatcher DomainMatcherGroup
	otherMatchers []matcherEntry
}

func NewMatcherGroup() *MatcherGroup {
	return &MatcherGroup{}
}

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

func (g *MatcherGroup) Size() uint32 {
	return g.count
}

type cacheEntry struct {
	timestamp time.Time
	result    uint32
}

type CachedMatcherGroup struct {
	sync.RWMutex
	group   *MatcherGroup
	cache   map[string]cacheEntry
	cleanup *task.Periodic
}

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

			expire := time.Now().Add(-1 * time.Second * 120)
			for p, e := range r.cache {
				if e.timestamp.Before(expire) {
					delete(r.cache, p)
				}
			}

			return nil
		},
	}
	return r
}

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
