package core

import "net"
import "sync"

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

func (d *syncDNSClient) Close() {
	d.RLock()
	defer d.RUnlock()

	if d.DNSClient != nil {
		d.DNSClient.Close()
	}
}

func (d *syncDNSClient) Set(client DNSClient) {
	if client == nil {
		return
	}

	d.Lock()
	defer d.Unlock()

	d.DNSClient = client
}
