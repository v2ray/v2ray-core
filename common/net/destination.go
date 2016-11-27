package net

import (
	"net"
)

// Destination represents a network destination including address and protocol (tcp / udp).
type Destination struct {
	Network Network
	Address Address
	Port    Port
}

func DestinationFromAddr(addr net.Addr) Destination {
	switch addr := addr.(type) {
	case *net.TCPAddr:
		return TCPDestination(IPAddress(addr.IP), Port(addr.Port))
	case *net.UDPAddr:
		return UDPDestination(IPAddress(addr.IP), Port(addr.Port))
	default:
		panic("Unknown address type.")
	}
}

// TCPDestination creates a TCP destination with given address
func TCPDestination(address Address, port Port) Destination {
	return Destination{
		Network: Network_TCP,
		Address: address,
		Port:    port,
	}
}

// UDPDestination creates a UDP destination with given address
func UDPDestination(address Address, port Port) Destination {
	return Destination{
		Network: Network_UDP,
		Address: address,
		Port:    port,
	}
}

func (v Destination) NetAddr() string {
	return v.Address.String() + ":" + v.Port.String()
}

func (v Destination) String() string {
	return v.Network.UrlPrefix() + ":" + v.NetAddr()
}

func (v *Endpoint) AsDestination() Destination {
	return Destination{
		Network: v.Network,
		Address: v.Address.AsAddress(),
		Port:    Port(v.Port),
	}
}
