package net

import (
	"github.com/v2ray/v2ray-core/common/log"
)

const (
	NetTCP = byte(0x01)
	NetUDP = byte(0x02)
)

type Destination struct {
	network byte
	address Address
}

func NewDestination(network byte, address Address) *Destination {
	return &Destination{
		network: network,
		address: address,
	}
}

func (dest *Destination) Network() string {
	switch dest.network {
	case NetTCP:
		return "tcp"
	case NetUDP:
		return "udp"
	default:
		log.Warning("Unknown network %d", dest.network)
		return "tcp"
	}
}

func (dest *Destination) NetworkByte() byte {
	return dest.network
}

func (dest *Destination) Address() Address {
	return dest.address
}

func (dest *Destination) String() string {
	return dest.address.String() + " (" + dest.Network() + ")"
}
