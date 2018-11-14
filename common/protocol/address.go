package protocol

import (
	"io"

	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/serial"
)

type AddressOption func(*option)

func PortThenAddress() AddressOption {
	return func(p *option) {
		p.portFirst = true
	}
}

func AddressFamilyByte(b byte, f net.AddressFamily) AddressOption {
	if b >= 16 {
		panic("address family byte too big")
	}
	return func(p *option) {
		p.addrTypeMap[b] = f
		p.addrByteMap[f] = b
	}
}

type AddressTypeParser func(byte) byte

func WithAddressTypeParser(atp AddressTypeParser) AddressOption {
	return func(p *option) {
		p.typeParser = atp
	}
}

type AddressSerializer interface {
	ReadAddressPort(buffer *buf.Buffer, input io.Reader) (net.Address, net.Port, error)

	WriteAddressPort(writer io.Writer, addr net.Address, port net.Port) error
}

const afInvalid = 255

type option struct {
	addrTypeMap [16]net.AddressFamily
	addrByteMap [16]byte
	portFirst   bool
	typeParser  AddressTypeParser
}

// NewAddressParser creates a new AddressParser
func NewAddressParser(options ...AddressOption) AddressSerializer {
	var o option
	for i := range o.addrByteMap {
		o.addrByteMap[i] = afInvalid
	}
	for i := range o.addrTypeMap {
		o.addrTypeMap[i] = net.AddressFamily(afInvalid)
	}
	for _, opt := range options {
		opt(&o)
	}

	ap := &addressParser{
		addrByteMap: o.addrByteMap,
		addrTypeMap: o.addrTypeMap,
	}

	if o.typeParser != nil {
		ap.typeParser = o.typeParser
	}

	if o.portFirst {
		return portFirstAddressParser{ap: ap}
	}

	return portLastAddressParser{ap: ap}
}

type portFirstAddressParser struct {
	ap *addressParser
}

func (p portFirstAddressParser) ReadAddressPort(buffer *buf.Buffer, input io.Reader) (net.Address, net.Port, error) {
	if buffer == nil {
		buffer = buf.New()
		defer buffer.Release()
	}

	port, err := readPort(buffer, input)
	if err != nil {
		return nil, 0, err
	}

	addr, err := p.ap.readAddress(buffer, input)
	if err != nil {
		return nil, 0, err
	}
	return addr, port, nil
}

func (p portFirstAddressParser) WriteAddressPort(writer io.Writer, addr net.Address, port net.Port) error {
	if err := writePort(writer, port); err != nil {
		return err
	}

	return p.ap.writeAddress(writer, addr)
}

type portLastAddressParser struct {
	ap *addressParser
}

func (p portLastAddressParser) ReadAddressPort(buffer *buf.Buffer, input io.Reader) (net.Address, net.Port, error) {
	if buffer == nil {
		buffer = buf.New()
		defer buffer.Release()
	}

	addr, err := p.ap.readAddress(buffer, input)
	if err != nil {
		return nil, 0, err
	}

	port, err := readPort(buffer, input)
	if err != nil {
		return nil, 0, err
	}

	return addr, port, nil
}

func (p portLastAddressParser) WriteAddressPort(writer io.Writer, addr net.Address, port net.Port) error {
	if err := p.ap.writeAddress(writer, addr); err != nil {
		return err
	}

	return writePort(writer, port)
}

func readPort(b *buf.Buffer, reader io.Reader) (net.Port, error) {
	if _, err := b.ReadFullFrom(reader, 2); err != nil {
		return 0, err
	}
	return net.PortFromBytes(b.BytesFrom(-2)), nil
}

func writePort(writer io.Writer, port net.Port) error {
	return common.Error2(serial.WriteUint16(writer, port.Value()))
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

type addressParser struct {
	addrTypeMap [16]net.AddressFamily
	addrByteMap [16]byte
	typeParser  AddressTypeParser
}

func (p *addressParser) readAddress(b *buf.Buffer, reader io.Reader) (net.Address, error) {
	if _, err := b.ReadFullFrom(reader, 1); err != nil {
		return nil, err
	}

	addrType := b.Byte(b.Len() - 1)
	if p.typeParser != nil {
		addrType = p.typeParser(addrType)
	}

	if addrType >= 16 {
		return nil, newError("unknown address type: ", addrType)
	}

	addrFamily := p.addrTypeMap[addrType]
	if addrFamily == net.AddressFamily(afInvalid) {
		return nil, newError("unknown address type: ", addrType)
	}

	switch addrFamily {
	case net.AddressFamilyIPv4:
		if _, err := b.ReadFullFrom(reader, 4); err != nil {
			return nil, err
		}
		return net.IPAddress(b.BytesFrom(-4)), nil
	case net.AddressFamilyIPv6:
		if _, err := b.ReadFullFrom(reader, 16); err != nil {
			return nil, err
		}
		return net.IPAddress(b.BytesFrom(-16)), nil
	case net.AddressFamilyDomain:
		if _, err := b.ReadFullFrom(reader, 1); err != nil {
			return nil, err
		}
		domainLength := int32(b.Byte(b.Len() - 1))
		if _, err := b.ReadFullFrom(reader, domainLength); err != nil {
			return nil, err
		}
		domain := string(b.BytesFrom(-domainLength))
		if maybeIPPrefix(domain[0]) {
			addr := net.ParseAddress(domain)
			if addr.Family().IsIPv4() || addr.Family().IsIPv6() {
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

func (p *addressParser) writeAddress(writer io.Writer, address net.Address) error {
	tb := p.addrByteMap[address.Family()]
	if tb == afInvalid {
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
	default:
		panic("Unknown family type.")
	}

	return nil
}
