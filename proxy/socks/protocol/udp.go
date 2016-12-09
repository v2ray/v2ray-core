package protocol

import (
	"fmt"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/errors"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/common/serial"
)

var (
	ErrorUnknownAddressType = errors.New("Unknown Address Type.")
)

type Socks5UDPRequest struct {
	Fragment byte
	Address  v2net.Address
	Port     v2net.Port
	Data     *buf.Buffer
}

func (request *Socks5UDPRequest) Destination() v2net.Destination {
	return v2net.UDPDestination(request.Address, request.Port)
}

func (request *Socks5UDPRequest) Write(buffer *buf.Buffer) {
	buffer.AppendBytes(0, 0, request.Fragment)
	switch request.Address.Family() {
	case v2net.AddressFamilyIPv4:
		buffer.AppendBytes(AddrTypeIPv4)
		buffer.Append(request.Address.IP())
	case v2net.AddressFamilyIPv6:
		buffer.AppendBytes(AddrTypeIPv6)
		buffer.Append(request.Address.IP())
	case v2net.AddressFamilyDomain:
		buffer.AppendBytes(AddrTypeDomain, byte(len(request.Address.Domain())))
		buffer.Append([]byte(request.Address.Domain()))
	}
	buffer.AppendSupplier(serial.WriteUint16(request.Port.Value()))
	buffer.Append(request.Data.Bytes())
}

func ReadUDPRequest(packet []byte) (*Socks5UDPRequest, error) {
	if len(packet) < 5 {
		return nil, errors.New("Socks|UDP: Insufficient length of packet.")
	}
	request := new(Socks5UDPRequest)

	// packet[0] and packet[1] are reserved
	request.Fragment = packet[2]

	addrType := packet[3]
	var dataBegin int

	switch addrType {
	case AddrTypeIPv4:
		if len(packet) < 10 {
			return nil, errors.New("Socks|UDP: Insufficient length of packet.")
		}
		ip := packet[4:8]
		request.Port = v2net.PortFromBytes(packet[8:10])
		request.Address = v2net.IPAddress(ip)
		dataBegin = 10
	case AddrTypeIPv6:
		if len(packet) < 22 {
			return nil, errors.New("Socks|UDP: Insufficient length of packet.")
		}
		ip := packet[4:20]
		request.Port = v2net.PortFromBytes(packet[20:22])
		request.Address = v2net.IPAddress(ip)
		dataBegin = 22
	case AddrTypeDomain:
		domainLength := int(packet[4])
		if len(packet) < 5+domainLength+2 {
			return nil, errors.New("Socks|UDP: Insufficient length of packet.")
		}
		domain := string(packet[5 : 5+domainLength])
		request.Port = v2net.PortFromBytes(packet[5+domainLength : 5+domainLength+2])
		request.Address = v2net.ParseAddress(domain)
		dataBegin = 5 + domainLength + 2
	default:
		return nil, fmt.Errorf("Socks|UDP: Unknown address type %d", addrType)
	}

	if len(packet) > dataBegin {
		b := buf.NewSmall()
		b.Append(packet[dataBegin:])
		request.Data = b
	}

	return request, nil
}
