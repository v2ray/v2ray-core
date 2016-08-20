package shadowsocks

import (
	"bytes"
	"io"

	"v2ray.com/core/common/alloc"
	"v2ray.com/core/common/log"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/proxy"
	"v2ray.com/core/transport"
)

const (
	AddrTypeIPv4   = 1
	AddrTypeIPv6   = 4
	AddrTypeDomain = 3
)

type Request struct {
	Address    v2net.Address
	Port       v2net.Port
	OTA        bool
	UDPPayload *alloc.Buffer
}

func (this *Request) Release() {
	this.Address = nil
	if this.UDPPayload != nil {
		this.UDPPayload.Release()
	}
}

func (this *Request) DetachUDPPayload() *alloc.Buffer {
	payload := this.UDPPayload
	this.UDPPayload = nil
	return payload
}

func ReadRequest(reader io.Reader, auth *Authenticator, udp bool) (*Request, error) {
	buffer := alloc.NewSmallBuffer()
	defer buffer.Release()

	_, err := io.ReadFull(reader, buffer.Value[:1])
	if err != nil {
		if err != io.EOF {
			log.Warning("Shadowsocks: Failed to read address type: ", err)
			return nil, transport.ErrCorruptedPacket
		}
		return nil, err
	}
	lenBuffer := 1

	request := new(Request)

	addrType := (buffer.Value[0] & 0x0F)
	if (buffer.Value[0] & 0x10) == 0x10 {
		request.OTA = true
	}
	switch addrType {
	case AddrTypeIPv4:
		_, err := io.ReadFull(reader, buffer.Value[lenBuffer:lenBuffer+4])
		if err != nil {
			log.Warning("Shadowsocks: Failed to read IPv4 address: ", err)
			return nil, transport.ErrCorruptedPacket
		}
		request.Address = v2net.IPAddress(buffer.Value[lenBuffer : lenBuffer+4])
		lenBuffer += 4
	case AddrTypeIPv6:
		_, err := io.ReadFull(reader, buffer.Value[lenBuffer:lenBuffer+16])
		if err != nil {
			log.Warning("Shadowsocks: Failed to read IPv6 address: ", err)
			return nil, transport.ErrCorruptedPacket
		}
		request.Address = v2net.IPAddress(buffer.Value[lenBuffer : lenBuffer+16])
		lenBuffer += 16
	case AddrTypeDomain:
		_, err := io.ReadFull(reader, buffer.Value[lenBuffer:lenBuffer+1])
		if err != nil {
			log.Warning("Shadowsocks: Failed to read domain lenth: ", err)
			return nil, transport.ErrCorruptedPacket
		}
		domainLength := int(buffer.Value[lenBuffer])
		lenBuffer++
		_, err = io.ReadFull(reader, buffer.Value[lenBuffer:lenBuffer+domainLength])
		if err != nil {
			log.Warning("Shadowsocks: Failed to read domain: ", err)
			return nil, transport.ErrCorruptedPacket
		}
		request.Address = v2net.DomainAddress(string(buffer.Value[lenBuffer : lenBuffer+domainLength]))
		lenBuffer += domainLength
	default:
		log.Warning("Shadowsocks: Unknown address type: ", addrType)
		return nil, transport.ErrCorruptedPacket
	}

	_, err = io.ReadFull(reader, buffer.Value[lenBuffer:lenBuffer+2])
	if err != nil {
		log.Warning("Shadowsocks: Failed to read port: ", err)
		return nil, transport.ErrCorruptedPacket
	}

	request.Port = v2net.PortFromBytes(buffer.Value[lenBuffer : lenBuffer+2])
	lenBuffer += 2

	var authBytes []byte

	if udp {
		nBytes, err := reader.Read(buffer.Value[lenBuffer:])
		if err != nil {
			log.Warning("Shadowsocks: Failed to read UDP payload: ", err)
			return nil, transport.ErrCorruptedPacket
		}
		buffer.Slice(0, lenBuffer+nBytes)
		if request.OTA {
			authBytes = buffer.Value[lenBuffer+nBytes-AuthSize:]
			request.UDPPayload = alloc.NewSmallBuffer().Clear().Append(buffer.Value[lenBuffer : lenBuffer+nBytes-AuthSize])
			lenBuffer = lenBuffer + nBytes - AuthSize
		} else {
			request.UDPPayload = alloc.NewSmallBuffer().Clear().Append(buffer.Value[lenBuffer:])
		}
	} else {
		if request.OTA {
			authBytes = buffer.Value[lenBuffer : lenBuffer+AuthSize]
			_, err = io.ReadFull(reader, authBytes)
			if err != nil {
				log.Warning("Shadowsocks: Failed to read OTA: ", err)
				return nil, transport.ErrCorruptedPacket
			}
		}
	}

	if request.OTA {
		actualAuth := auth.Authenticate(nil, buffer.Value[0:lenBuffer])
		if !bytes.Equal(actualAuth, authBytes) {
			log.Warning("Shadowsocks: Invalid OTA.")
			return nil, proxy.ErrInvalidAuthentication
		}
	}

	return request, nil
}
