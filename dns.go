package core

import (
	"sync"

	"v2ray.com/core/common"
	"v2ray.com/core/common/net"
	"v2ray.com/core/features/dns"
)

type syncDNSClient struct {
	sync.RWMutex
	dns.Client
}

func (d *syncDNSClient) Type() interface{} {
	return dns.ClientType()
}

func (d *syncDNSClient) LookupIP(host string) ([]net.IP, error) {
	d.RLock()
	defer d.RUnlock()

	if d.Client == nil {
		return net.LookupIP(host)
	}

	return d.Client.LookupIP(host)
}

func (d *syncDNSClient) Start() error {
	d.RLock()
	defer d.RUnlock()

	if d.Client == nil {
		return nil
	}

	return d.Client.Start()
}

func (d *syncDNSClient) Close() error {
	d.RLock()
	defer d.RUnlock()

	return common.Close(d.Client)
}

func (d *syncDNSClient) Set(client dns.Client) {
	if client == nil {
		return
	}

	d.Lock()
	defer d.Unlock()

	common.Close(d.Client) // nolint: errcheck
	d.Client = client
}
