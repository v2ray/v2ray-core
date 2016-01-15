package json

import (
	"encoding/json"
	"errors"
	"net"
	"regexp"
	"strings"

	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/common/serial"
)

type DomainMatcher interface {
	Match(domain string) bool
}

type PlainDomainMatcher struct {
	pattern string
}

func NewPlainDomainMatcher(pattern string) *PlainDomainMatcher {
	return &PlainDomainMatcher{
		pattern: strings.ToLower(pattern),
	}
}

func (this *PlainDomainMatcher) Match(domain string) bool {
	return strings.Contains(strings.ToLower(domain), this.pattern)
}

type RegexpDomainMatcher struct {
	pattern *regexp.Regexp
}

func NewRegexpDomainMatcher(pattern string) (*RegexpDomainMatcher, error) {
	r, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	return &RegexpDomainMatcher{
		pattern: r,
	}, nil
}

func (this *RegexpDomainMatcher) Match(domain string) bool {
	return this.pattern.MatchString(strings.ToLower(domain))
}

type FieldRule struct {
	Rule
	Domain  []DomainMatcher
	IP      []*net.IPNet
	Port    *v2net.PortRange
	Network *v2net.NetworkList
}

func (this *FieldRule) Apply(dest v2net.Destination) bool {
	address := dest.Address()
	if len(this.Domain) > 0 {
		if !address.IsDomain() {
			return false
		}
		foundMatch := false
		for _, domain := range this.Domain {
			if domain.Match(address.Domain()) {
				foundMatch = true
				break
			}
		}
		if !foundMatch {
			return false
		}
	}

	if this.IP != nil && len(this.IP) > 0 {
		if !(address.IsIPv4() || address.IsIPv6()) {
			return false
		}
		foundMatch := false
		for _, ipnet := range this.IP {
			if ipnet.Contains(address.IP()) {
				foundMatch = true
				break
			}
		}
		if !foundMatch {
			return false
		}
	}

	if this.Port != nil {
		port := dest.Port()
		if port.Value() < this.Port.From.Value() || port.Value() > this.Port.To.Value() {
			return false
		}
	}

	if this.Network != nil {
		if !this.Network.HasNetwork(v2net.Network(dest.Network())) {
			return false
		}
	}

	return true
}

func (this *FieldRule) UnmarshalJSON(data []byte) error {
	type RawFieldRule struct {
		Rule
		Domain  *serial.StringLiteralList `json:"domain"`
		IP      *serial.StringLiteralList `json:"ip"`
		Port    *v2net.PortRange          `json:"port"`
		Network *v2net.NetworkList        `json:"network"`
	}
	rawFieldRule := RawFieldRule{}
	err := json.Unmarshal(data, &rawFieldRule)
	if err != nil {
		return err
	}
	this.Type = rawFieldRule.Type
	this.OutboundTag = rawFieldRule.OutboundTag

	hasField := false
	if rawFieldRule.Domain != nil && rawFieldRule.Domain.Len() > 0 {
		this.Domain = make([]DomainMatcher, rawFieldRule.Domain.Len())
		for idx, rawDomain := range *(rawFieldRule.Domain) {
			var matcher DomainMatcher
			if strings.HasPrefix(rawDomain.String(), "regexp:") {
				rawMatcher, err := NewRegexpDomainMatcher(rawDomain.String()[7:])
				if err != nil {
					return err
				}
				matcher = rawMatcher
			} else {
				matcher = NewPlainDomainMatcher(rawDomain.String())
			}
			this.Domain[idx] = matcher
		}
		hasField = true
	}

	if rawFieldRule.IP != nil && rawFieldRule.IP.Len() > 0 {
		this.IP = make([]*net.IPNet, 0, rawFieldRule.IP.Len())
		for _, ipStr := range *(rawFieldRule.IP) {
			_, ipNet, err := net.ParseCIDR(ipStr.String())
			if err != nil {
				return errors.New("Invalid IP range in router rule: " + err.Error())
			}
			this.IP = append(this.IP, ipNet)
		}
		hasField = true
	}
	if rawFieldRule.Port != nil {
		this.Port = rawFieldRule.Port
		hasField = true
	}
	if rawFieldRule.Network != nil {
		this.Network = rawFieldRule.Network
		hasField = true
	}
	if !hasField {
		return errors.New("This rule has no effective fields.")
	}
	return nil
}
