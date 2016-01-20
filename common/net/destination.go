package net

// Destination represents a network destination including address and protocol (tcp / udp).
type Destination interface {
	Network() string  // Protocol of communication (tcp / udp)
	Address() Address // Address of destination
	Port() Port
	String() string // String representation of the destination
	NetAddr() string
	Equals(Destination) bool

	IsTCP() bool // True if destination is reachable via TCP
	IsUDP() bool // True if destination is reachable via UDP
}

// TCPDestination creates a TCP destination with given address
func TCPDestination(address Address, port Port) Destination {
	return &tcpDestination{address: address, port: port}
}

// UDPDestination creates a UDP destination with given address
func UDPDestination(address Address, port Port) Destination {
	return &udpDestination{address: address, port: port}
}

type tcpDestination struct {
	address Address
	port    Port
}

func (dest *tcpDestination) Network() string {
	return "tcp"
}

func (dest *tcpDestination) Address() Address {
	return dest.address
}

func (dest *tcpDestination) NetAddr() string {
	return dest.address.String() + ":" + dest.port.String()
}

func (dest *tcpDestination) String() string {
	return "tcp:" + dest.NetAddr()
}

func (dest *tcpDestination) IsTCP() bool {
	return true
}

func (dest *tcpDestination) IsUDP() bool {
	return false
}

func (dest *tcpDestination) Port() Port {
	return dest.port
}

func (dest *tcpDestination) Equals(another Destination) bool {
	if !another.IsTCP() {
		return false
	}
	return dest.Port() == another.Port() && dest.Address().Equals(another.Address())
}

type udpDestination struct {
	address Address
	port    Port
}

func (dest *udpDestination) Network() string {
	return "udp"
}

func (dest *udpDestination) Address() Address {
	return dest.address
}

func (dest *udpDestination) NetAddr() string {
	return dest.address.String() + ":" + dest.port.String()
}

func (dest *udpDestination) String() string {
	return "udp:" + dest.NetAddr()
}

func (dest *udpDestination) IsTCP() bool {
	return false
}

func (dest *udpDestination) IsUDP() bool {
	return true
}

func (dest *udpDestination) Port() Port {
	return dest.port
}

func (dest *udpDestination) Equals(another Destination) bool {
	if !another.IsUDP() {
		return false
	}
	return dest.Port() == another.Port() && dest.Address().Equals(another.Address())
}
