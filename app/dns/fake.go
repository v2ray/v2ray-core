package dns

import (
	"bufio"
	"container/list"
	"os"
	"strconv"
	"strings"

	"v2ray.com/core/common/net"
	"v2ray.com/core/common/strmatcher"
)

func max(a uint32, b uint32) uint32 {
	if a > b {
		return a
	}
	return b
}

func getIPSum(i net.IP) uint32 {
	return (uint32(i[0]) << 24) | (uint32(i[1]) << 16) | (uint32(i[2]) << 8) | uint32(i[3])
}

func getAddress(i uint32) net.Address {
	ip := make([]byte, 4)
	ip[0] = byte(i >> 24)
	ip[1] = byte(i >> 16)
	ip[2] = byte(i >> 8)
	ip[3] = byte(i)
	return net.IPAddress(ip)
}

var matcher *strmatcher.OrMatcher

// Prefix of fake ip index
var prefix uint32
var upperLimit int

func isFakeIP(i uint32) bool {
	return (i & prefix) == prefix
}

type fakeIPMapper interface {
	clear()
	getAddress(string) []net.Address
	getDomain(uint32) string
	load(map[uint32]string)
}

type noneMapper struct {
	domainMapper  map[string][]net.Address
	addressMapper map[uint32]string
	next          uint32
}

func (n *noneMapper) clear() {
	n.domainMapper = make(map[string][]net.Address)
	n.addressMapper = make(map[uint32]string)
	n.next = 0
}

func (n *noneMapper) getAddress(domain string) []net.Address {
	if res := n.domainMapper[domain]; res != nil {
		return res
	}
	ipIndex := prefix | n.next
	if n.next >= uint32(upperLimit-1) {
		return nil
	}
	go saveFakeIP(n.next, domain)
	n.next++
	ret := []net.Address{getAddress(ipIndex)}
	n.addressMapper[ipIndex] = domain
	n.domainMapper[domain] = ret
	return ret
}

func (n *noneMapper) getDomain(ipIndex uint32) string {
	if len(n.addressMapper) == 0 {
		return ""
	}

	return n.addressMapper[ipIndex]
}

func (n *noneMapper) load(rules map[uint32]string) {
	for ip, domain := range rules {
		if !isFakeIP(ip) {
			continue
		}
		address := []net.Address{getAddress(ip)}
		n.addressMapper[ip] = domain
		n.domainMapper[domain] = address
		n.next = max(n.next, ip)
	}
}

type oldestMapper struct {
	noneMapper
}

func (n *oldestMapper) getAddress(domain string) []net.Address {
	if n.next >= uint32(upperLimit-1) {
		n.next = 0
	} else {
		n.next++
	}
	return n.noneMapper.getAddress(domain)
}

type lruNode struct {
	domain  *domainNode
	address *addressNode
}
type domainNode struct {
	domain string
	lru    *list.Element
}
type addressNode struct {
	address []net.Address
	lru     *list.Element
}
type lruMapper struct {
	domainMapper  map[string]*addressNode
	addressMapper map[uint32]*domainNode
	lru           *list.List
	next          uint32
}

func (n *lruMapper) goNextEmpty() {
	for {
		_, isOk := n.addressMapper[n.next]
		if !isOk {
			return
		}
		n.next++
	}
}

func (n *lruMapper) clear() {
	n.domainMapper = make(map[string]*addressNode)
	n.addressMapper = make(map[uint32]*domainNode)
	n.lru = list.New()
	n.next = prefix
}

func (n *lruMapper) getAddress(domain string) []net.Address {
	res := n.domainMapper[domain]
	if res != nil {
		n.lru.MoveBefore(res.lru, n.lru.Front())
		return res.address
	}
	var lru *list.Element
	if len(n.addressMapper) >= upperLimit {
		lru = n.lru.Back()
		go saveFakeAddress(lru.Value.(*lruNode).address.address, domain)
		n.lru.MoveBefore(lru, n.lru.Front())
		delete(n.domainMapper, lru.Value.(*lruNode).domain.domain)
	} else {
		res = new(addressNode)
		dom := new(domainNode)
		lru = n.lru.PushFront(&lruNode{
			domain:  dom,
			address: res,
		})
		n.goNextEmpty()
		go saveFakeIP(n.next, domain)
		res.address = []net.Address{getAddress(n.next)}
		res.lru = lru
		n.addressMapper[n.next] = dom
		dom.lru = lru
	}
	lru.Value.(*lruNode).domain.domain = domain
	n.domainMapper[domain] = lru.Value.(*lruNode).address
	return lru.Value.(*lruNode).address.address
}

func (n *lruMapper) getDomain(ipIndex uint32) string {
	if len(n.addressMapper) == 0 {
		return ""
	}
	res := n.addressMapper[ipIndex]
	if res == nil {
		return ""
	}
	n.lru.MoveBefore(res.lru, n.lru.Front())
	return res.domain
}

func (n *lruMapper) load(rules map[uint32]string) {
	for ip, domain := range rules {
		if !isFakeIP(ip) {
			continue
		}
		res := new(addressNode)
		dom := new(domainNode)
		lru := n.lru.PushFront(&lruNode{
			domain:  dom,
			address: res,
		})
		res.address = []net.Address{getAddress(ip)}
		res.lru = lru
		n.addressMapper[ip] = dom
		dom.lru = lru
		lru.Value.(*lruNode).domain.domain = domain
		n.domainMapper[domain] = lru.Value.(*lruNode).address
	}
}

var fakeIP fakeIPMapper
var saver *os.File

func saveFakeAddress(address []net.Address, domain string) {
	if saver == nil {
		return
	}
	saveFakeIP(getIPSum(address[0].IP()), domain)
}

func saveFakeIP(ip uint32, domain string) {
	if saver == nil {
		return
	}
	saver.WriteString(strconv.Itoa(int(ip)) + " " + domain + "\n")
}

func loadFakeIP() (map[uint32]string, error) {
	ipToDomain := make(map[uint32]string)
	scanner := bufio.NewScanner(saver)
	for scanner.Scan() {
		rule := strings.Split(scanner.Text(), " ")
		ip, err := strconv.Atoi(rule[0])
		if err != nil {
			return nil, err
		}
		ipToDomain[uint32(ip)] = rule[1]
	}
	return ipToDomain, nil
}

func prepareFakeFile(path string) (map[uint32]string, error) {
	var err error
	saver, err = os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0600)
	if err != nil {
		return nil, newError("failed to open fake ip file: ", path).Base(err).AtWarning()
	}
	rules, err := loadFakeIP()
	if err != nil {
		return nil, newError("fake ip file corrupted: ", path).Base(err).AtWarning()
	}
	saver.Close()
	saver, err = os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return nil, newError("failed to open fake ip file: ", path).Base(err).AtWarning()
	}
	// Delete duplicate rules
	for ip, domain := range rules {
		saveFakeIP(ip, domain)
	}
	return rules, nil
}

// InitFakeIPServer initializes matcher for domain name checking
func InitFakeIPServer(fake *Config_Fake, externalRules map[string][]string) error {
	if fake != nil {
		if fake.FakeRules == nil {
			return newError("no rules for fake ip").AtWarning()
		}
		nd := strings.Split(fake.FakeNet, "/")
		mask, err := strconv.Atoi(nd[1])
		if err != nil {
			return newError("failed to parse fakeNet: ", fake.FakeNet).Base(err).AtWarning()
		}
		upperLimit = 1 << (32 - mask)
		prefix = getIPSum(net.ParseAddress(nd[0]).IP()) & uint32(^(upperLimit - 1))
		switch fake.Regeneration {
		case Config_Fake_LRU:
			fakeIP = new(lruMapper)
		case Config_Fake_Oldest:
			fakeIP = new(oldestMapper)
		case Config_Fake_None:
			fakeIP = new(noneMapper)
		}
		ResetFakeIPServer()
		for _, pattern := range fake.FakeRules {
			if err := matcher.ParsePattern(pattern, externalRules); err != nil {
				newError("failed to parse pattern: ", pattern).Base(err).AtWarning().WriteToLog()
			}
		}
		if fake.Path != "" {
			rules, err := prepareFakeFile(fake.Path)
			if err != nil {
				return err
			}
			fakeIP.load(rules)
		}
	}
	return nil
}

// GetFakeIPForDomain checks if we should response with a fake IP for a domain name
func GetFakeIPForDomain(domain string) []net.Address {
	if matcher == nil || !matcher.Match(domain) {
		return nil
	}

	return fakeIP.getAddress(domain)
}

// GetDomainForFakeIP checks if a IP is a fake IP and return its corresponding domain name
func GetDomainForFakeIP(ip net.Address) string {
	if fakeIP == nil || !ip.Family().IsIP() {
		return ""
	}
	sum := getIPSum(ip.IP())
	if isFakeIP(sum) {
		return fakeIP.getDomain(sum)
	}
	return ""
}

// ResetFakeIPServer is for testing only
func ResetFakeIPServer() {
	matcher = strmatcher.NewOrMatcher()
	fakeIP.clear()
}
