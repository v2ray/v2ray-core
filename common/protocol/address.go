package protocol

import (
	"io"

	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/net"
)

type AddressOption func(*AddressParser)

func PortThenAddress() AddressOption {
	return func(p *AddressParser) {
		p.portFirst = true
	}
}

func AddressFamilyByte(b byte, f net.AddressFamily) AddressOption {
	return func(p *AddressParser) {
		p.addrTypeMap[b] = f
		p.addrByteMap[f] = b
	}
}

type AddressTypeParser func(byte) byte

func WithAddressTypeParser(atp AddressTypeParser) AddressOption {
	return func(p *AddressParser) {
		p.typeParser = atp
	}
}

type AddressParser struct {
	addrTypeMap map[byte]net.AddressFamily
	addrByteMap map[net.AddressFamily]byte
	portFirst   bool
	typeParser  AddressTypeParser
}

func NewAddressParser(options ...AddressOption) *AddressParser {
	p := &AddressParser{
		addrTypeMap: make(map[byte]net.AddressFamily, 8),
		addrByteMap: make(map[net.AddressFamily]byte, 8),
	}
	for _, opt := range options {
		opt(p)
	}
	return p
}

func (p *AddressParser) readPort(b *buf.Buffer, reader io.Reader) (net.Port, error) {
	if err := b.AppendSupplier(buf.ReadFullFrom(reader, 2)); err != nil {
		return 0, err
	}
	return net.PortFromBytes(b.BytesFrom(-2)), nil
}

func maybeIPPrefix(b byte) bool {
	return b == '[' || (b >= '0' && b <= '9')
}

func isValidDomain(d string) bool {
	for _, c := range d {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '-' || c == '.' || c == '_') {
			return false
		}
	}
	return true
}

func (p *AddressParser) readAddress(b *buf.Buffer, reader io.Reader) (net.Address, error) {
	if err := b.AppendSupplier(buf.ReadFullFrom(reader, 1)); err != nil {
		return nil, err
	}

	addrType := b.Byte(b.Len() - 1)
	if p.typeParser != nil {
		addrType = p.typeParser(addrType)
	}

	addrFamily, valid := p.addrTypeMap[addrType]
	if !valid {
		return nil, newError("unknown address type: ", addrType)
	}

	switch addrFamily {
	case net.AddressFamilyIPv4:
		if err := b.AppendSupplier(buf.ReadFullFrom(reader, 4)); err != nil {
			return nil, err
		}
		return net.IPAddress(b.BytesFrom(-4)), nil
	case net.AddressFamilyIPv6:
		if err := b.AppendSupplier(buf.ReadFullFrom(reader, 16)); err != nil {
			return nil, err
		}
		return net.IPAddress(b.BytesFrom(-16)), nil
	case net.AddressFamilyDomain:
		if err := b.AppendSupplier(buf.ReadFullFrom(reader, 1)); err != nil {
			return nil, err
		}
		domainLength := int32(b.Byte(b.Len() - 1))
		if err := b.AppendSupplier(buf.ReadFullFrom(reader, domainLength)); err != nil {
			return nil, err
		}
		domain := string(b.BytesFrom(-domainLength))
		if maybeIPPrefix(domain[0]) {
			addr := net.ParseAddress(domain)
			if addr.Family().IsIPv4() || addrFamily.IsIPv6() {
				return addr, nil
			}
		}
		if !isValidDomain(domain) {
			return nil, newError("invalid domain name: ", domain)
		}
		return net.DomainAddress(domain), nil
	default:
		panic("impossible case")
	}
}

func (p *AddressParser) ReadAddressPort(buffer *buf.Buffer, input io.Reader) (net.Address, net.Port, error) {
	if buffer == nil {
		buffer = buf.New()
		defer buffer.Release()
	}

	if p.portFirst {
		port, err := p.readPort(buffer, input)
		if err != nil {
			return nil, 0, err
		}
		addr, err := p.readAddress(buffer, input)
		if err != nil {
			return nil, 0, err
		}
		return addr, port, nil
	}

	addr, err := p.readAddress(buffer, input)
	if err != nil {
		return nil, 0, err
	}

	port, err := p.readPort(buffer, input)
	if err != nil {
		return nil, 0, err
	}

	return addr, port, nil
}

func (p *AddressParser) writePort(writer io.Writer, port net.Port) error {
	_, err := writer.Write(port.Bytes(nil))
	return err
}

func (p *AddressParser) writeAddress(writer io.Writer, address net.Address) error {
	tb, valid := p.addrByteMap[address.Family()]
	if !valid {
		return newError("unknown address family", address.Family())
	}

	switch address.Family() {
	case net.AddressFamilyIPv4, net.AddressFamilyIPv6:
		if _, err := writer.Write([]byte{tb}); err != nil {
			return err
		}
		if _, err := writer.Write(address.IP()); err != nil {
			return err
		}
	case net.AddressFamilyDomain:
		domain := address.Domain()
		if isDomainTooLong(domain) {
			return newError("Super long domain is not supported: ", domain)
		}
		if _, err := writer.Write([]byte{tb, byte(len(domain))}); err != nil {
			return err
		}
		if _, err := writer.Write([]byte(domain)); err != nil {
			return err
		}
	}
	return nil
}

func (p *AddressParser) WriteAddressPort(writer io.Writer, addr net.Address, port net.Port) error {
	if p.portFirst {
		if err := p.writePort(writer, port); err != nil {
			return err
		}
		if err := p.writeAddress(writer, addr); err != nil {
			return err
		}
		return nil
	}

	if err := p.writeAddress(writer, addr); err != nil {
		return err
	}
	if err := p.writePort(writer, port); err != nil {
		return err
	}
	return nil
}
