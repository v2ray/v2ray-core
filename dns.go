package core

import (
	"net"
	"sync"

	"v2ray.com/core/common"
)

// DNSClient is a V2Ray feature for querying DNS information.
type DNSClient interface {
	Feature
	LookupIP(host string) ([]net.IP, error)
}

type syncDNSClient struct {
	sync.RWMutex
	DNSClient
}

func (d *syncDNSClient) LookupIP(host string) ([]net.IP, error) {
	d.RLock()
	defer d.RUnlock()

	if d.DNSClient == nil {
		return net.LookupIP(host)
	}

	return d.DNSClient.LookupIP(host)
}

func (d *syncDNSClient) Start() error {
	d.RLock()
	defer d.RUnlock()

	if d.DNSClient == nil {
		return nil
	}

	return d.DNSClient.Start()
}

func (d *syncDNSClient) Close() error {
	d.RLock()
	defer d.RUnlock()

	return common.Close(d.DNSClient)
}

func (d *syncDNSClient) Set(client DNSClient) {
	if client == nil {
		return
	}

	d.Lock()
	defer d.Unlock()

	d.DNSClient = client
}
