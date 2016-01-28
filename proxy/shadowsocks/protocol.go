package shadowsocks

import (
	"io"

	"github.com/v2ray/v2ray-core/common/alloc"
	"github.com/v2ray/v2ray-core/common/log"
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/transport"
)

const (
	AddrTypeIPv4   = 1
	AddrTypeIPv6   = 4
	AddrTypeDomain = 3
)

type Request struct {
	Address v2net.Address
	Port    v2net.Port
}

func ReadRequest(reader io.Reader) (*Request, error) {
	buffer := alloc.NewSmallBuffer()
	defer buffer.Release()

	_, err := v2net.ReadAllBytes(reader, buffer.Value[:1])
	if err != nil {
		log.Error("Shadowsocks: Failed to read address type: ", err)
		return nil, transport.CorruptedPacket
	}

	request := new(Request)

	addrType := buffer.Value[0]
	switch addrType {
	case AddrTypeIPv4:
		_, err := v2net.ReadAllBytes(reader, buffer.Value[:4])
		if err != nil {
			log.Error("Shadowsocks: Failed to read IPv4 address: ", err)
			return nil, transport.CorruptedPacket
		}
		request.Address = v2net.IPAddress(buffer.Value[:4])
	case AddrTypeIPv6:
		_, err := v2net.ReadAllBytes(reader, buffer.Value[:16])
		if err != nil {
			log.Error("Shadowsocks: Failed to read IPv6 address: ", err)
			return nil, transport.CorruptedPacket
		}
		request.Address = v2net.IPAddress(buffer.Value[:16])
	case AddrTypeDomain:
		_, err := v2net.ReadAllBytes(reader, buffer.Value[:1])
		if err != nil {
			log.Error("Shadowsocks: Failed to read domain lenth: ", err)
			return nil, transport.CorruptedPacket
		}
		domainLength := int(buffer.Value[0])
		_, err = v2net.ReadAllBytes(reader, buffer.Value[:domainLength])
		if err != nil {
			log.Error("Shadowsocks: Failed to read domain: ", err)
			return nil, transport.CorruptedPacket
		}
		request.Address = v2net.DomainAddress(string(buffer.Value[:domainLength]))
	default:
		log.Error("Shadowsocks: Unknown address type: ", addrType)
		return nil, transport.CorruptedPacket
	}

	_, err = v2net.ReadAllBytes(reader, buffer.Value[:2])
	if err != nil {
		log.Error("Shadowsocks: Failed to read port: ", err)
		return nil, transport.CorruptedPacket
	}

	request.Port = v2net.PortFromBytes(buffer.Value[:2])
	return request, nil
}
