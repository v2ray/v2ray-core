package strmatcher

import "strings"

func breakDomain(domain string) []string {
	return strings.Split(domain, ".")
}

type node struct {
	value uint32
	sub   map[string]*node
}

type DomainMatcherGroup struct {
	root *node
}

func (g *DomainMatcherGroup) Add(domain string, value uint32) {
	if g.root == nil {
		g.root = &node{
			sub: make(map[string]*node),
		}
	}

	current := g.root
	parts := breakDomain(domain)
	for i := len(parts) - 1; i >= 0; i-- {
		part := parts[i]
		next := current.sub[part]
		if next == nil {
			next = &node{sub: make(map[string]*node)}
			current.sub[part] = next
		}
		current = next
	}

	current.value = value
}

func (g *DomainMatcherGroup) Match(domain string) uint32 {
	current := g.root
	parts := breakDomain(domain)
	for i := len(parts) - 1; i >= 0; i-- {
		part := parts[i]
		next := current.sub[part]
		if next == nil {
			break
		}
		current = next
	}
	return current.value
}
