package net

// Destination represents a network destination including address and protocol (tcp / udp).
type Destination interface {
	Network() string  // Protocol of communication (tcp / udp)
	Address() Address // Address of destination
	String() string   // String representation of the destination

	IsTCP() bool // True if destination is reachable via TCP
	IsUDP() bool // True if destination is reachable via UDP
}

// NewTCPDestination creates a TCP destination with given address
func NewTCPDestination(address Address) Destination {
	return TCPDestination{address: address}
}

// NewUDPDestination creates a UDP destination with given address
func NewUDPDestination(address Address) Destination {
	return UDPDestination{address: address}
}

type TCPDestination struct {
	address Address
}

func (dest TCPDestination) Network() string {
	return "tcp"
}

func (dest TCPDestination) Address() Address {
	return dest.address
}

func (dest TCPDestination) String() string {
	return "tcp:" + dest.address.String()
}

func (dest TCPDestination) IsTCP() bool {
	return true
}

func (dest TCPDestination) IsUDP() bool {
	return false
}

type UDPDestination struct {
	address Address
}

func (dest UDPDestination) Network() string {
	return "udp"
}

func (dest UDPDestination) Address() Address {
	return dest.address
}

func (dest UDPDestination) String() string {
	return "udp:" + dest.address.String()
}

func (dest UDPDestination) IsTCP() bool {
	return false
}

func (dest UDPDestination) IsUDP() bool {
	return true
}
