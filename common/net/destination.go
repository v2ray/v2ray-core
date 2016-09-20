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

func (this Destination) NetAddr() string {
	return this.Address.String() + ":" + this.Port.String()
}

func (this Destination) String() string {
	return this.Network.UrlPrefix() + ":" + this.NetAddr()
}

func (this Destination) Equals(another Destination) bool {
	return this.Network == another.Network && this.Port == another.Port && this.Address.Equals(another.Address)
}

func (this *DestinationPB) AsDestination() Destination {
	return Destination{
		Network: this.Network,
		Address: this.Address.AsAddress(),
		Port:    Port(this.Port),
	}
}
