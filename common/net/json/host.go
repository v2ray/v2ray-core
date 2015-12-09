package json

import (
	"encoding/json"
	"net"
)

type Host struct {
	domain string
	ip     net.IP
}

func NewIPHost(ip net.IP) *Host {
	return &Host{
		ip: ip,
	}
}

func NewDomainHost(domain string) *Host {
	return &Host{
		domain: domain,
	}
}

func (this *Host) UnmarshalJSON(data []byte) error {
	var rawStr string
	if err := json.Unmarshal(data, &rawStr); err != nil {
		return err
	}
	ip := net.ParseIP(rawStr)
	if ip != nil {
		this.ip = ip
	} else {
		this.domain = rawStr
	}
	return nil
}

func (this *Host) IsIP() bool {
	return this.ip != nil
}

func (this *Host) IsDomain() bool {
	return !this.IsIP()
}

func (this *Host) IP() net.IP {
	return this.ip
}

func (this *Host) Domain() string {
	return this.domain
}
