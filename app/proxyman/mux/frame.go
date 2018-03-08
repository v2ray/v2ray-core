package mux

import (
	"v2ray.com/core/common/bitmask"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/protocol"
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

var addrParser = protocol.NewAddressParser(
	protocol.AddressFamilyByte(byte(protocol.AddressTypeIPv4), net.AddressFamilyIPv4),
	protocol.AddressFamilyByte(byte(protocol.AddressTypeDomain), net.AddressFamilyDomain),
	protocol.AddressFamilyByte(byte(protocol.AddressTypeIPv6), net.AddressFamilyIPv6),
	protocol.PortThenAddress(),
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

func (f FrameMetadata) WriteTo(b *buf.Buffer) error {
	lenBytes := b.Bytes()
	b.AppendBytes(0x00, 0x00)

	len0 := b.Len()
	if err := b.AppendSupplier(serial.WriteUint16(f.SessionID)); err != nil {
		return err
	}

	b.AppendBytes(byte(f.SessionStatus), byte(f.Option))

	if f.SessionStatus == SessionStatusNew {
		switch f.Target.Network {
		case net.Network_TCP:
			b.AppendBytes(byte(TargetNetworkTCP))
		case net.Network_UDP:
			b.AppendBytes(byte(TargetNetworkUDP))
		}

		if err := addrParser.WriteAddressPort(b, f.Target.Address, f.Target.Port); err != nil {
			return err
		}
	}

	len1 := b.Len()
	serial.Uint16ToBytes(uint16(len1-len0), lenBytes)
	return nil
}

func ReadFrameFrom(b *buf.Buffer) (*FrameMetadata, error) {
	if b.Len() < 4 {
		return nil, newError("insufficient buffer: ", b.Len())
	}

	f := &FrameMetadata{
		SessionID:     serial.BytesToUint16(b.BytesTo(2)),
		SessionStatus: SessionStatus(b.Byte(2)),
		Option:        bitmask.Byte(b.Byte(3)),
	}

	if f.SessionStatus == SessionStatusNew {
		network := TargetNetwork(b.Byte(4))
		b.SliceFrom(5)

		addr, port, err := addrParser.ReadAddressPort(nil, b)
		if err != nil {
			return nil, newError("failed to parse address and port").Base(err)
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
