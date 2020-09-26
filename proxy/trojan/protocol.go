package trojan

import (
	"encoding/binary"
	"io"

	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/protocol"
)

var (
	crlf = []byte{'\r', '\n'}

	addrParser = protocol.NewAddressParser(
		protocol.AddressFamilyByte(0x01, net.AddressFamilyIPv4),   // nolint: gomnd
		protocol.AddressFamilyByte(0x04, net.AddressFamilyIPv6),   // nolint: gomnd
		protocol.AddressFamilyByte(0x03, net.AddressFamilyDomain), // nolint: gomnd
	)
)

const (
	maxLength = 8192

	commandTCP byte = 1
	commandUDP byte = 3
)

// ConnWriter is TCP Connection Writer Wrapper for trojan protocol
type ConnWriter struct {
	io.Writer
	Target     net.Destination
	Account    *MemoryAccount
	headerSent bool
}

// Write implements io.Writer
func (c *ConnWriter) Write(p []byte) (n int, err error) {
	if !c.headerSent {
		if err := c.writeHeader(); err != nil {
			return 0, newError("failed to write request header").Base(err)
		}
	}

	return c.Writer.Write(p)
}

// WriteMultiBuffer implements buf.Writer
func (c *ConnWriter) WriteMultiBuffer(mb buf.MultiBuffer) error {
	defer buf.ReleaseMulti(mb)

	for _, b := range mb {
		if !b.IsEmpty() {
			if _, err := c.Write(b.Bytes()); err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *ConnWriter) writeHeader() error {
	buffer := buf.StackNew()
	defer buffer.Release()

	command := commandTCP
	if c.Target.Network == net.Network_UDP {
		command = commandUDP
	}

	if _, err := buffer.Write(c.Account.Key); err != nil {
		return err
	}
	if _, err := buffer.Write(crlf); err != nil {
		return err
	}
	if err := buffer.WriteByte(command); err != nil {
		return err
	}
	if err := addrParser.WriteAddressPort(&buffer, c.Target.Address, c.Target.Port); err != nil {
		return err
	}
	if _, err := buffer.Write(crlf); err != nil {
		return err
	}

	_, err := c.Writer.Write(buffer.Bytes())
	if err == nil {
		c.headerSent = true
	}

	return err
}

// PacketWriter UDP Connection Writer Wrapper for trojan protocol
type PacketWriter struct {
	io.Writer
	Target net.Destination
}

// WriteMultiBuffer implements buf.Writer
func (w *PacketWriter) WriteMultiBuffer(mb buf.MultiBuffer) error {
	b := make([]byte, maxLength)
	for !mb.IsEmpty() {
		var length int
		mb, length = buf.SplitBytes(mb, b)
		if _, err := w.writePacket(b[:length], w.Target); err != nil {
			buf.ReleaseMulti(mb)
			return err
		}
	}

	return nil
}

// WriteMultiBufferWithMetadata writes udp packet with destination specified
func (w *PacketWriter) WriteMultiBufferWithMetadata(mb buf.MultiBuffer, dest net.Destination) error {
	b := make([]byte, maxLength)
	for !mb.IsEmpty() {
		var length int
		mb, length = buf.SplitBytes(mb, b)
		if _, err := w.writePacket(b[:length], dest); err != nil {
			buf.ReleaseMulti(mb)
			return err
		}
	}

	return nil
}

func (w *PacketWriter) writePacket(payload []byte, dest net.Destination) (int, error) { // nolint: unparam
	buffer := buf.StackNew()
	defer buffer.Release()

	length := len(payload)
	lengthBuf := [2]byte{}
	binary.BigEndian.PutUint16(lengthBuf[:], uint16(length))
	if err := addrParser.WriteAddressPort(&buffer, dest.Address, dest.Port); err != nil {
		return 0, err
	}
	if _, err := buffer.Write(lengthBuf[:]); err != nil {
		return 0, err
	}
	if _, err := buffer.Write(crlf); err != nil {
		return 0, err
	}
	if _, err := buffer.Write(payload); err != nil {
		return 0, err
	}
	_, err := w.Write(buffer.Bytes())
	if err != nil {
		return 0, err
	}

	return length, nil
}

// ConnReader is TCP Connection Reader Wrapper for trojan protocol
type ConnReader struct {
	io.Reader
	Target       net.Destination
	headerParsed bool
}

// ParseHeader parses the trojan protocol header
func (c *ConnReader) ParseHeader() error {
	var crlf [2]byte
	var command [1]byte
	var hash [56]byte
	if _, err := io.ReadFull(c.Reader, hash[:]); err != nil {
		return newError("failed to read user hash").Base(err)
	}

	if _, err := io.ReadFull(c.Reader, crlf[:]); err != nil {
		return newError("failed to read crlf").Base(err)
	}

	if _, err := io.ReadFull(c.Reader, command[:]); err != nil {
		return newError("failed to read command").Base(err)
	}

	network := net.Network_TCP
	if command[0] == commandUDP {
		network = net.Network_UDP
	}

	addr, port, err := addrParser.ReadAddressPort(nil, c.Reader)
	if err != nil {
		return newError("failed to read address and port").Base(err)
	}
	c.Target = net.Destination{Network: network, Address: addr, Port: port}

	if _, err := io.ReadFull(c.Reader, crlf[:]); err != nil {
		return newError("failed to read crlf").Base(err)
	}

	c.headerParsed = true
	return nil
}

// Read implements io.Reader
func (c *ConnReader) Read(p []byte) (int, error) {
	if !c.headerParsed {
		if err := c.ParseHeader(); err != nil {
			return 0, err
		}
	}

	return c.Reader.Read(p)
}

// ReadMultiBuffer implements buf.Reader
func (c *ConnReader) ReadMultiBuffer() (buf.MultiBuffer, error) {
	b := buf.New()
	_, err := b.ReadFrom(c)
	return buf.MultiBuffer{b}, err
}

// PacketPayload combines udp payload and destination
type PacketPayload struct {
	Target net.Destination
	Buffer buf.MultiBuffer
}

// PacketReader is UDP Connection Reader Wrapper for trojan protocol
type PacketReader struct {
	io.Reader
}

// ReadMultiBuffer implements buf.Reader
func (r *PacketReader) ReadMultiBuffer() (buf.MultiBuffer, error) {
	p, err := r.ReadMultiBufferWithMetadata()
	if p != nil {
		return p.Buffer, err
	}
	return nil, err
}

// ReadMultiBufferWithMetadata reads udp packet with destination
func (r *PacketReader) ReadMultiBufferWithMetadata() (*PacketPayload, error) {
	addr, port, err := addrParser.ReadAddressPort(nil, r)
	if err != nil {
		return nil, newError("failed to read address and port").Base(err)
	}

	var lengthBuf [2]byte
	if _, err := io.ReadFull(r, lengthBuf[:]); err != nil {
		return nil, newError("failed to read payload length").Base(err)
	}

	remain := int(binary.BigEndian.Uint16(lengthBuf[:]))
	if remain > maxLength {
		return nil, newError("oversize payload")
	}

	var crlf [2]byte
	if _, err := io.ReadFull(r, crlf[:]); err != nil {
		return nil, newError("failed to read crlf").Base(err)
	}

	dest := net.UDPDestination(addr, port)
	var mb buf.MultiBuffer
	for remain > 0 {
		length := buf.Size
		if remain < length {
			length = remain
		}

		b := buf.New()
		mb = append(mb, b)
		n, err := b.ReadFullFrom(r, int32(length))
		if err != nil {
			buf.ReleaseMulti(mb)
			return nil, newError("failed to read payload").Base(err)
		}

		remain -= int(n)
	}

	return &PacketPayload{Target: dest, Buffer: mb}, nil
}
