package protocol

import (
	"io"

	"v2ray.com/core/common/task"

	"v2ray.com/core/common"
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

// AddressParser is a utility for reading and writer addresses.
type AddressParser struct {
	addrTypeMap map[byte]net.AddressFamily
	addrByteMap map[net.AddressFamily]byte
	portFirst   bool
	typeParser  AddressTypeParser
}

// NewAddressParser creates a new AddressParser
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

// ReadAddressPort reads address and port from the given input.
func (p *AddressParser) ReadAddressPort(buffer *buf.Buffer, input io.Reader) (net.Address, net.Port, error) {
	if buffer == nil {
		buffer = buf.New()
		defer buffer.Release()
	}

	var addr net.Address
	var port net.Port

	pTask := func() error {
		lp, err := p.readPort(buffer, input)
		if err != nil {
			return err
		}
		port = lp
		return nil
	}

	aTask := func() error {
		a, err := p.readAddress(buffer, input)
		if err != nil {
			return err
		}
		addr = a
		return nil
	}

	var err error

	if p.portFirst {
		err = task.Run(task.Sequential(pTask, aTask))()
	} else {
		err = task.Run(task.Sequential(aTask, pTask))()
	}

	if err != nil {
		return nil, 0, err
	}

	return addr, port, nil
}

func (p *AddressParser) writePort(writer io.Writer, port net.Port) error {
	return common.Error2(writer.Write(port.Bytes(nil)))
}

func (p *AddressParser) writeAddress(writer io.Writer, address net.Address) error {
	tb, valid := p.addrByteMap[address.Family()]
	if !valid {
		return newError("unknown address family", address.Family())
	}

	switch address.Family() {
	case net.AddressFamilyIPv4, net.AddressFamilyIPv6:
		return task.Run(task.Sequential(func() error {
			return common.Error2(writer.Write([]byte{tb}))
		}, func() error {
			return common.Error2(writer.Write(address.IP()))
		}))()
	case net.AddressFamilyDomain:
		domain := address.Domain()
		if isDomainTooLong(domain) {
			return newError("Super long domain is not supported: ", domain)
		}
		return task.Run(task.Sequential(func() error {
			return common.Error2(writer.Write([]byte{tb, byte(len(domain))}))
		}, func() error {
			return common.Error2(writer.Write([]byte(domain)))
		}))()
	default:
		panic("Unknown family type.")
	}
}

// WriteAddressPort writes address and port into the given writer.
func (p *AddressParser) WriteAddressPort(writer io.Writer, addr net.Address, port net.Port) error {
	pTask := func() error {
		return p.writePort(writer, port)
	}
	aTask := func() error {
		return p.writeAddress(writer, addr)
	}

	if p.portFirst {
		return task.Run(task.Sequential(pTask, aTask))()
	}

	return task.Run(task.Sequential(aTask, pTask))()
}
