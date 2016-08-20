package net

import (
	"strings"

	"v2ray.com/core/common/collect"
)

const (
	// TCPNetwork represents the TCP network.
	TCPNetwork = Network("tcp")

	// UDPNetwork represents the UDP network.
	UDPNetwork = Network("udp")

	// KCPNetwork represents the KCP network.
	KCPNetwork = Network("kcp")

	// WSNetwork represents the Websocket over HTTP network.
	WSNetwork = Network("ws")
)

// Network represents a communication network on internet.
type Network string

func (this Network) AsList() *NetworkList {
	list := NetworkList([]Network{this})
	return &list
}

func (this Network) String() string {
	return string(this)
}

// NetworkList is a list of Networks.
type NetworkList []Network

// NewNetworkList construsts a NetWorklist from the given StringListeralList.
func NewNetworkList(networks collect.StringList) NetworkList {
	list := NetworkList(make([]Network, networks.Len()))
	for idx, network := range networks {
		list[idx] = Network(strings.ToLower(strings.TrimSpace(network)))
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
