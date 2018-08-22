package strmatcher

import "strings"

func breakDomain(domain string) []string {
	return strings.Split(domain, ".")
}

type node struct {
	value uint32
	sub   map[string]*node
}

// DomainMatcherGroup is a IndexMatcher for a large set of Domain matchers.
// Visible for testing only.
type DomainMatcherGroup struct {
	root *node
}

func (g *DomainMatcherGroup) Add(domain string, value uint32) {
	if g.root == nil {
		g.root = new(node)
	}

	current := g.root
	parts := breakDomain(domain)
	for i := len(parts) - 1; i >= 0; i-- {
		part := parts[i]
		if current.sub == nil {
			current.sub = make(map[string]*node)
		}
		next := current.sub[part]
		if next == nil {
			next = new(node)
			current.sub[part] = next
		}
		current = next
	}

	current.value = value
}

func (g *DomainMatcherGroup) addMatcher(m domainMatcher, value uint32) {
	g.Add(string(m), value)
}

func (g *DomainMatcherGroup) Match(domain string) uint32 {
	current := g.root
	if current == nil {
		return 0
	}

	parts := breakDomain(domain)
	for i := len(parts) - 1; i >= 0; i-- {
		part := parts[i]
		if current.sub == nil {
			break
		}
		next := current.sub[part]
		if next == nil {
			break
		}
		current = next
	}
	return current.value
}
