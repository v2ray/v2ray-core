package mux

import (
	"v2ray.com/core/common/bitmask"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/serial"
)

type SessionStatus byte

const (
	SessionStatusNew       SessionStatus = 0x01
	SessionStatusKeep      SessionStatus = 0x02
	SessionStatusEnd       SessionStatus = 0x03
	SessionStatusKeepAlive SessionStatus = 0x04
)

const (
	OptionData bitmask.Byte = 0x01
)

type TargetNetwork byte

const (
	TargetNetworkTCP TargetNetwork = 0x01
	TargetNetworkUDP TargetNetwork = 0x02
)

type AddressType byte

const (
	AddressTypeIPv4   AddressType = 0x01
	AddressTypeDomain AddressType = 0x02
	AddressTypeIPv6   AddressType = 0x03
)

/*
Frame format
2 bytes - length
2 bytes - session id
1 bytes - status
1 bytes - option

1 byte - network
2 bytes - port
n bytes - address

*/

type FrameMetadata struct {
	Target        net.Destination
	SessionID     uint16
	Option        bitmask.Byte
	SessionStatus SessionStatus
}

func (f FrameMetadata) AsSupplier() buf.Supplier {
	return func(b []byte) (int, error) {
		lengthBytes := b
		b = serial.Uint16ToBytes(uint16(0), b[:0]) // place holder for length

		b = serial.Uint16ToBytes(f.SessionID, b)
		b = append(b, byte(f.SessionStatus), byte(f.Option))
		length := 4

		if f.SessionStatus == SessionStatusNew {
			switch f.Target.Network {
			case net.Network_TCP:
				b = append(b, byte(TargetNetworkTCP))
			case net.Network_UDP:
				b = append(b, byte(TargetNetworkUDP))
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
				b = append(b, byte(AddressTypeDomain), byte(nDomain))
				b = append(b, addr.Domain()...)
				length += nDomain + 2
			}
		}

		serial.Uint16ToBytes(uint16(length), lengthBytes[:0])
		return length + 2, nil
	}
}

func ReadFrameFrom(b []byte) (*FrameMetadata, error) {
	if len(b) < 4 {
		return nil, newError("insufficient buffer: ", len(b))
	}

	f := &FrameMetadata{
		SessionID:     serial.BytesToUint16(b[:2]),
		SessionStatus: SessionStatus(b[2]),
		Option:        bitmask.Byte(b[3]),
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
		default:
			return nil, newError("unknown address type: ", addrType)
		}
		switch network {
		case TargetNetworkTCP:
			f.Target = net.TCPDestination(addr, port)
		case TargetNetworkUDP:
			f.Target = net.UDPDestination(addr, port)
		default:
			return nil, newError("unknown network type: ", network)
		}
	}

	return f, nil
}
