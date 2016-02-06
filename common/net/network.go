package net

import (
	"github.com/v2ray/v2ray-core/common/serial"
)

const (
	// TCPNetwork represents the TCP network.
	TCPNetwork = Network("tcp")

	// UDPNetwork represents the UDP network.
	UDPNetwork = Network("udp")
)

// Network represents a communication network on internet.
type Network serial.StringLiteral

func (this Network) AsList() *NetworkList {
	list := NetworkList([]Network{this})
	return &list
}

// NetworkList is a list of Networks.
type NetworkList []Network

// NewNetworkList construsts a NetWorklist from the given StringListeralList.
func NewNetworkList(networks serial.StringLiteralList) NetworkList {
	list := NetworkList(make([]Network, networks.Len()))
	for idx, network := range networks {
		list[idx] = Network(network.TrimSpace().ToLower())
	}
	return list
}

// HashNetwork returns true if the given network is in this NetworkList.
func (this *NetworkList) HasNetwork(network Network) bool {
	for _, value := range *this {
		if string(value) == string(network) {
			return true
		}
	}
	return false
}
