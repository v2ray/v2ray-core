package net

const (
	TCPNetwork = Network("tcp")
	UDPNetwork = Network("udp")
)

type Network string

type NetworkList interface {
	HasNetwork(Network) bool
}
