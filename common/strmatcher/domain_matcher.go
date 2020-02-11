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
		if current.value > 0 {
			// if current node is already a match, it is not necessary to match further.
			return
		}

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
	current.sub = nil // shortcut sub nodes as current node is a match.
}

func (g *DomainMatcherGroup) addMatcher(m domainMatcher, value uint32) {
	g.Add(string(m), value)
}

func (g *DomainMatcherGroup) Match(domain string) uint32 {
	if domain == "" {
		return 0
	}

	current := g.root
	if current == nil {
		return 0
	}

	nextPart := func(idx int) int {
		for i := idx - 1; i >= 0; i-- {
			if domain[i] == '.' {
				return i
			}
		}
		return -1
	}

	idx := len(domain)
	for {
		if idx == -1 || current.sub == nil {
			break
		}

		nidx := nextPart(idx)
		part := domain[nidx+1 : idx]
		next := current.sub[part]
		if next == nil {
			break
		}
		current = next
		idx = nidx
	}
	return current.value
}

type domainNode struct {
	sub map[string]*domainNode
}

type domainGroupMatcher struct {
	root *domainNode
}

func (g *domainGroupMatcher) Add(domain string) {
	if g.root == nil {
		g.root = new(domainNode)
	}
	current := g.root
	parts := breakDomain(domain)
	for i := len(parts) - 1; i >= 0; i-- {
		// if current node is already a match, it is not necessary to match further.
		if current.sub != nil && len(current.sub) == 0 {
			return
		}
		if current.sub == nil {
			current.sub = make(map[string]*domainNode)
		}
		part := parts[i]
		next := current.sub[part]
		if next == nil {
			next = new(domainNode)
			current.sub[part] = next
		}
		current = next
	}

	current.sub = make(map[string]*domainNode) // shortcut sub nodes as current node is a match.
}

func (g *domainGroupMatcher) addMatcher(m domainMatcher) {
	g.Add(string(m))
}

func (g *domainGroupMatcher) Match(domain string) bool {
	if domain == "" {
		return false
	}

	current := g.root
	if current == nil {
		return false
	}

	nextPart := func(idx int) int {
		for i := idx - 1; i >= 0; i-- {
			if domain[i] == '.' {
				return i
			}
		}
		return -1
	}

	idx := len(domain)
	for {
		if len(current.sub) == 0 {
			return true
		}
		if idx == -1 {
			return false
		}

		nidx := nextPart(idx)
		part := domain[nidx+1 : idx]
		next := current.sub[part]
		if next == nil {
			return false
		}
		current = next
		idx = nidx
	}
}
