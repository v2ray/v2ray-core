package net

import (
	"strings"
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

func (v Network) AsList() *NetworkList {
	return &NetworkList{
		Network: []Network{v},
	}
}

func (v Network) SystemString() string {
	switch v {
	case Network_TCP, Network_RawTCP:
		return "tcp"
	case Network_UDP, Network_KCP:
		return "udp"
	default:
		return "unknown"
	}
}

func (v Network) URLPrefix() string {
	switch v {
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

// HashNetwork returns true if the given network is in v NetworkList.
func (v NetworkList) HasNetwork(network Network) bool {
	for _, value := range v.Network {
		if string(value) == string(network) {
			return true
		}
	}
	return false
}

func (v NetworkList) Get(idx int) Network {
	return v.Network[idx]
}
