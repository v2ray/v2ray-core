package command

import (
	"errors"
	"io"

	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/transport"
)

func init() {
	RegisterResponseCommand(2, func() Command { return new(CacheDns) })
}

const (
	typeIPv4 byte = 1
	typeIPv6 byte = 2
)

var (
	ErrDomainAddress = errors.New("Unexpected domain address")
)

// Size: 1 byte type + 4 or 16 byte IP addr
type CacheDns struct {
	Address v2net.Address
}

func (this *CacheDns) Marshal(writer io.Writer) (int, error) {
	if this.Address.IsIPv4() {
		writer.Write([]byte{typeIPv4})
		writer.Write(this.Address.IP())
		return 5, nil
	}

	if this.Address.IsIPv6() {
		writer.Write([]byte{typeIPv6})
		writer.Write(this.Address.IP())
		return 17, nil
	}

	return 0, ErrDomainAddress
}

func (this *CacheDns) Unmarshal(data []byte) error {
	if len(data) == 0 {
		return transport.CorruptedPacket
	}
	typeIP := data[0]
	data = data[1:]

	if typeIP == typeIPv4 {
		if len(data) < 4 {
			return transport.CorruptedPacket
		}
		this.Address = v2net.IPAddress(data[0:4])
		return nil
	}

	if typeIP == typeIPv6 {
		if len(data) < 16 {
			return transport.CorruptedPacket
		}
		this.Address = v2net.IPAddress(data[0:16])
		return nil
	}

	return transport.CorruptedPacket
}
