package net

import (
	"strings"

	"v2ray.com/core/common/collect"
)

func ParseNetwork(nwStr string) Network {
	if network, found := Network_value[nwStr]; found {
		return Network(network)
	}
	switch strings.ToLower(nwStr) {
	case "tcp":
		return Network_TCP
	case "udp":
		return Network_UDP
	case "kcp":
		return Network_KCP
	case "ws":
		return Network_WebSocket
	default:
		return Network_Unknown
	}
}

func (this Network) AsList() *NetworkList {
	return &NetworkList{
		Network: []Network{this},
	}
}

func (this Network) SystemString() string {
	switch this {
	case Network_TCP, Network_RawTCP:
		return "tcp"
	case Network_UDP, Network_KCP:
		return "udp"
	default:
		return "unknown"
	}
}

func (this Network) UrlPrefix() string {
	switch this {
	case Network_TCP, Network_RawTCP:
		return "tcp"
	case Network_UDP:
		return "udp"
	case Network_KCP:
		return "kcp"
	case Network_WebSocket:
		return "ws"
	default:
		return "unknown"
	}
}

// NewNetworkList construsts a NetWorklist from the given StringListeralList.
func NewNetworkList(networks collect.StringList) *NetworkList {
	list := &NetworkList{
		Network: make([]Network, networks.Len()),
	}
	for idx, network := range networks {
		list.Network[idx] = ParseNetwork(network)
	}
	return list
}

// HashNetwork returns true if the given network is in this NetworkList.
func (this *NetworkList) HasNetwork(network Network) bool {
	for _, value := range this.Network {
		if string(value) == string(network) {
			return true
		}
	}
	return false
}
