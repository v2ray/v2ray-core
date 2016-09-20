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
	list := NetworkList([]Network{this})
	return &list
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

// NetworkList is a list of Networks.
type NetworkList []Network

// NewNetworkList construsts a NetWorklist from the given StringListeralList.
func NewNetworkList(networks collect.StringList) NetworkList {
	list := NetworkList(make([]Network, networks.Len()))
	for idx, network := range networks {
		list[idx] = ParseNetwork(network)
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
