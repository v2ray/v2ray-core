package net

import (
	"strings"
)

const (
	TCPNetwork = Network("tcp")
	UDPNetwork = Network("udp")
)

type Network string

type NetworkList []Network

func NewNetworkList(networks []string) NetworkList {
	list := NetworkList(make([]Network, len(networks)))
	for idx, network := range networks {
		list[idx] = Network(strings.ToLower(strings.TrimSpace(network)))
	}
	return list
}

func (this *NetworkList) HasNetwork(network Network) bool {
	for _, value := range *this {
		if string(value) == string(network) {
			return true
		}
	}
	return false
}
