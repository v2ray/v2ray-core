package protocol

import (
	"errors"
	"net"

	"github.com/v2ray/v2ray-core/common/alloc"
	"github.com/v2ray/v2ray-core/common/log"
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/transport"
)

var (
	ErrorUnknownAddressType = errors.New("Unknown Address Type.")
)

type Socks5UDPRequest struct {
	Fragment byte
	Address  v2net.Address
	Port     v2net.Port
	Data     *alloc.Buffer
}

func (request *Socks5UDPRequest) Destination() v2net.Destination {
	return v2net.UDPDestination(request.Address, request.Port)
}

func (request *Socks5UDPRequest) Write(buffer *alloc.Buffer) {
	buffer.AppendBytes(0, 0, request.Fragment)
	switch {
	case request.Address.IsIPv4():
		buffer.AppendBytes(AddrTypeIPv4).Append(request.Address.IP())
	case request.Address.IsIPv6():
		buffer.AppendBytes(AddrTypeIPv6).Append(request.Address.IP())
	case request.Address.IsDomain():
		buffer.AppendBytes(AddrTypeDomain, byte(len(request.Address.Domain()))).Append([]byte(request.Address.Domain()))
	}
	buffer.Append(request.Port.Bytes())
	buffer.Append(request.Data.Value)
}

func ReadUDPRequest(packet []byte) (*Socks5UDPRequest, error) {
	if len(packet) < 5 {
		return nil, transport.CorruptedPacket
	}
	request := new(Socks5UDPRequest)

	// packet[0] and packet[1] are reserved
	request.Fragment = packet[2]

	addrType := packet[3]
	var dataBegin int

	switch addrType {
	case AddrTypeIPv4:
		if len(packet) < 10 {
			return nil, transport.CorruptedPacket
		}
		ip := packet[4:8]
		request.Port = v2net.PortFromBytes(packet[8:10])
		request.Address = v2net.IPAddress(ip)
		dataBegin = 10
	case AddrTypeIPv6:
		if len(packet) < 22 {
			return nil, transport.CorruptedPacket
		}
		ip := packet[4:20]
		request.Port = v2net.PortFromBytes(packet[20:22])
		request.Address = v2net.IPAddress(ip)
		dataBegin = 22
	case AddrTypeDomain:
		domainLength := int(packet[4])
		if len(packet) < 5+domainLength+2 {
			return nil, transport.CorruptedPacket
		}
		domain := string(packet[5 : 5+domainLength])
		request.Port = v2net.PortFromBytes(packet[5+domainLength : 5+domainLength+2])
		maybeIP := net.ParseIP(domain)
		if maybeIP != nil {
			request.Address = v2net.IPAddress(maybeIP)
		} else {
			request.Address = v2net.DomainAddress(domain)
		}
		dataBegin = 5 + domainLength + 2
	default:
		log.Warning("Unknown address type ", addrType)
		return nil, ErrorUnknownAddressType
	}

	if len(packet) > dataBegin {
		request.Data = alloc.NewBuffer().Clear().Append(packet[dataBegin:])
	}

	return request, nil
}
