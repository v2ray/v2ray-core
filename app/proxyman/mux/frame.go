package mux

import (
	"errors"

	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/serial"
)

type SessionStatus byte

const (
	SessionStatusNew  SessionStatus = 0x01
	SessionStatusKeep SessionStatus = 0x02
	SessionStatusEnd  SessionStatus = 0x03
)

type TargetNetwork byte

const (
	TargetNetworkTCP TargetNetwork = 0x01
	TargetnetworkUDP TargetNetwork = 0x02
)

type AddressType byte

const (
	AddressTypeIPv4   AddressType = 0x01
	AddressTypeDomain AddressType = 0x02
	AddressTypeIPv6   AddressType = 0x03
)

type FrameMetadata struct {
	SessionId     uint16
	SessionStatus SessionStatus
	Target        net.Destination
}

func (f FrameMetadata) AsSupplier() buf.Supplier {
	return func(b []byte) (int, error) {
		b = serial.Uint16ToBytes(uint16(0), b) // place holder for length

		b = serial.Uint16ToBytes(f.SessionId, b)
		b = append(b, byte(f.SessionStatus), 0 /* reserved */)
		length := 4

		if f.SessionStatus == SessionStatusNew {
			switch f.Target.Network {
			case net.Network_TCP:
				b = append(b, byte(TargetNetworkTCP))
			case net.Network_UDP:
				b = append(b, byte(TargetnetworkUDP))
			}
			length++

			b = serial.Uint16ToBytes(f.Target.Port.Value(), b)
			length += 2

			addr := f.Target.Address
			switch addr.Family() {
			case net.AddressFamilyIPv4:
				b = append(b, byte(AddressTypeIPv4))
				b = append(b, addr.IP()...)
				length += 5
			case net.AddressFamilyIPv6:
				b = append(b, byte(AddressTypeIPv6))
				b = append(b, addr.IP()...)
				length += 17
			case net.AddressFamilyDomain:
				nDomain := len(addr.Domain())
				b = append(b, byte(nDomain))
				b = append(b, addr.Domain()...)
				length += nDomain + 1
			}
		}
		return length + 2, nil
	}
}

func ReadFrameFrom(b []byte) (*FrameMetadata, error) {
	if len(b) < 4 {
		return nil, errors.New("Proxyman|Mux: Insufficient buffer.")
	}

	f := &FrameMetadata{
		SessionId:     serial.BytesToUint16(b[:2]),
		SessionStatus: SessionStatus(b[2]),
	}

	b = b[4:]

	if f.SessionStatus == SessionStatusNew {
		network := TargetNetwork(b[0])
		port := net.PortFromBytes(b[1:3])
		addrType := AddressType(b[3])
		b = b[4:]

		var addr net.Address
		switch addrType {
		case AddressTypeIPv4:
			addr = net.IPAddress(b[0:4])
			b = b[4:]
		case AddressTypeIPv6:
			addr = net.IPAddress(b[0:16])
			b = b[16:]
		case AddressTypeDomain:
			nDomain := int(b[0])
			addr = net.DomainAddress(string(b[1 : 1+nDomain]))
			b = b[nDomain+1:]
		}
		switch network {
		case TargetNetworkTCP:
			f.Target = net.TCPDestination(addr, port)
		case TargetnetworkUDP:
			f.Target = net.UDPDestination(addr, port)
		}
	}

	return f, nil
}
