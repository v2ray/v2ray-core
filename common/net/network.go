package net

import (
	"github.com/v2ray/v2ray-core/common/serial"
)

const (
	TCPNetwork = Network("tcp")
	UDPNetwork = Network("udp")
)

type Network serial.StringLiteral

type NetworkList []Network

func NewNetworkList(networks serial.StringLiteralList) NetworkList {
	list := NetworkList(make([]Network, networks.Len()))
	for idx, network := range networks {
		list[idx] = Network(network.TrimSpace().ToLower())
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
