package net

import (
	"net"
	"sync/atomic"
	"unsafe"
)

// IPResolver is the interface to resolve host name to IPs.
type IPResolver interface {
	LookupIP(host string) ([]net.IP, error)
}

type systemIPResolver int

func (s systemIPResolver) LookupIP(host string) ([]net.IP, error) {
	return net.LookupIP(host)
}

const (
	systemIPResolverInstance = systemIPResolver(0)
)

// SystemIPResolver returns an IPResolver that resolves IP through underlying system.
func SystemIPResolver() IPResolver {
	return systemIPResolverInstance
}

var (
	ipResolver unsafe.Pointer
)

func LookupIP(host string) ([]net.IP, error) {
	r := (*IPResolver)(atomic.LoadPointer(&ipResolver))
	return (*r).LookupIP(host)
}

func RegisterIPResolver(resolver IPResolver) {
	atomic.StorePointer(&ipResolver, unsafe.Pointer(&resolver))
}

func init() {
	RegisterIPResolver(systemIPResolverInstance)
}
