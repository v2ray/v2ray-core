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
	"geoip:":   "i",
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

func compressPattern(pattern string, external map[string][]string, def string) (string, error) {
	for prefix, cmd := range prefixMapper {
		if strings.HasPrefix(pattern, prefix) {
			newPattern := cmd + pattern[len(prefix):]
			if newPattern[0] == 'e' {
				if err := loadExternalRules(newPattern[1:], external); err != nil {
					return "", err
				}
			}
			return newPattern, nil
		}
	}
	// If no prefix, use specified
	return def + pattern, nil
}
