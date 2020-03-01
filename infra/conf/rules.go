package conf

import (
	"strings"

	"v2ray.com/core/app/router"
)

var prefixMapper = map[string]string{
	"domain:":  "d",
	"regexp:":  "r",
	"keyword:": "k",
	"full:":    "f",
	"geosite:": "egeosite.dat:",
	"ext:":     "e",
	"not:":     "!",
	"and:":     "&",
}

var typeMapper = map[router.Domain_Type]string{
	router.Domain_Full:   "f",
	router.Domain_Domain: "d",
	router.Domain_Plain:  "k",
	router.Domain_Regex:  "r",
}

func loadExternalRules(pattern string, external map[string][]string) error {
	// Loaded rules
	if external[pattern] != nil {
		return nil
	}

	kv := strings.Split(pattern, ":")
	if len(kv) != 2 {
		return newError("invalid external resource: ", pattern)
	}
	filename, country := kv[0], kv[1]
	domains, err := loadGeositeWithAttr(filename, country)
	if err != nil {
		return newError("invalid external settings from ", filename, ": ", pattern).Base(err)
	}
	rule := make([]string, len(domains))
	index := 0
	for _, d := range domains {
		rule[index] = typeMapper[d.Type] + d.Value
		index++
	}

	external[pattern] = rule

	return nil
}

// In the nextPattern, parsePattern records the end position of the
// first pattern in the pattern it parsed.
type patternParser struct {
	unparsedNumber int
	nextPattern    int
	defaultType    string
	external       map[string][]string
}

func (p *patternParser) parsePattern(pattern string) (string, error) {
	for prefix, cmd := range prefixMapper {
		if !strings.HasPrefix(pattern, prefix) {
			continue
		}
		newPattern := cmd + pattern[len(prefix):]
		length := len(newPattern)
		// For the matchers which have not child matcher
		p.nextPattern = length
		if p.unparsedNumber != 0 && cmd != "&" && cmd != "!" {
			pos := strings.IndexByte(newPattern, ' ')
			if pos != -1 {
				p.nextPattern = pos
			}
		}
		arg := newPattern[1:p.nextPattern]
		switch newPattern[0] {
		case 'e':
			if err := loadExternalRules(arg, p.external); err != nil {
				return "", err
			}
		case '!':
			subPattern, err := p.parsePattern(arg)
			if err != nil {
				return "", err
			}
			// The sub pattern is one character shorter than the current string
			p.nextPattern++
			newPattern = "!" + subPattern
		case '&':
			p.unparsedNumber++
			partA, err := p.parsePattern(arg)
			if err != nil {
				return "", err
			}
			lenA := p.nextPattern
			// The part after p.nextPattern haven't been parsed yet
			newPattern = "&" + partA[:lenA]
			partB, err := p.parsePattern(partA[lenA+1:])
			if err != nil {
				return "", err
			}
			p.nextPattern = lenA + 2 + p.nextPattern
			p.unparsedNumber--
			newPattern += " " + partB
		}
		return newPattern, nil
	}
	// If no prefix, use specified
	return p.defaultType + pattern, nil
}

func compressPattern(pattern string, external map[string][]string, def string) (string, error) {
	p := &patternParser{
		unparsedNumber: 0,
		nextPattern:    0,
		defaultType:    def,
		external:       external,
	}
	return p.parsePattern(pattern)
}
