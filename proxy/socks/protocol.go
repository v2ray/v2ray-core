package socks

import (
	"encoding/binary"
	"io"

	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/protocol"
)

const (
	socks5Version = 0x05
	socks4Version = 0x04

	cmdTCPConnect    = 0x01
	cmdTCPBind       = 0x02
	cmdUDPPort       = 0x03
	cmdTorResolve    = 0xF0
	cmdTorResolvePTR = 0xF1

	socks4RequestGranted  = 90
	socks4RequestRejected = 91

	authNotRequired = 0x00
	//authGssAPI           = 0x01
	authPassword         = 0x02
	authNoMatchingMethod = 0xFF

	statusSuccess       = 0x00
	statusCmdNotSupport = 0x07
)

var addrParser = protocol.NewAddressParser(
	protocol.AddressFamilyByte(0x01, net.AddressFamilyIPv4),
	protocol.AddressFamilyByte(0x04, net.AddressFamilyIPv6),
	protocol.AddressFamilyByte(0x03, net.AddressFamilyDomain),
)

type ServerSession struct {
	config *ServerConfig
	port   net.Port
}

func (s *ServerSession) handshake4(cmd byte, reader io.Reader, writer io.Writer) (*protocol.RequestHeader, error) {
	if s.config.AuthType == AuthType_PASSWORD {
		writeSocks4Response(writer, socks4RequestRejected, net.AnyIP, net.Port(0)) // nolint: errcheck
		return nil, newError("socks 4 is not allowed when auth is required.")
	}

	var port net.Port
	var address net.Address

	{
		buffer := buf.StackNew()
		if _, err := buffer.ReadFullFrom(reader, 6); err != nil {
			buffer.Release()
			return nil, newError("insufficient header").Base(err)
		}
		port = net.PortFromBytes(buffer.BytesRange(0, 2))
		address = net.IPAddress(buffer.BytesRange(2, 6))
		buffer.Release()
	}

	if _, err := ReadUntilNull(reader); /* user id */ err != nil {
		return nil, err
	}
	if address.IP()[0] == 0x00 {
		domain, err := ReadUntilNull(reader)
		if err != nil {
			return nil, newError("failed to read domain for socks 4a").Base(err)
		}
		address = net.DomainAddress(domain)
	}

	switch cmd {
	case cmdTCPConnect:
		request := &protocol.RequestHeader{
			Command: protocol.RequestCommandTCP,
			Address: address,
			Port:    port,
			Version: socks4Version,
		}
		if err := writeSocks4Response(writer, socks4RequestGranted, net.AnyIP, net.Port(0)); err != nil {
			return nil, err
		}
		return request, nil
	default:
		writeSocks4Response(writer, socks4RequestRejected, net.AnyIP, net.Port(0)) // nolint: errcheck
		return nil, newError("unsupported command: ", cmd)
	}
}

func (s *ServerSession) auth5(nMethod byte, reader io.Reader, writer io.Writer) error {
	buffer := buf.StackNew()
	defer buffer.Release()

	if _, err := buffer.ReadFullFrom(reader, int32(nMethod)); err != nil {
		return newError("failed to read auth methods").Base(err)
	}

	var expectedAuth byte = authNotRequired
	if s.config.AuthType == AuthType_PASSWORD {
		expectedAuth = authPassword
	}

	if !hasAuthMethod(expectedAuth, buffer.BytesRange(0, int32(nMethod))) {
		writeSocks5AuthenticationResponse(writer, socks5Version, authNoMatchingMethod) // nolint: errcheck
		return newError("no matching auth method")
	}

	if err := writeSocks5AuthenticationResponse(writer, socks5Version, expectedAuth); err != nil {
		return newError("failed to write auth response").Base(err)
	}

	if expectedAuth == authPassword {
		username, password, err := ReadUsernamePassword(reader)
		if err != nil {
			return newError("failed to read username and password for authentication").Base(err)
		}

		if !s.config.HasAccount(username, password) {
			writeSocks5AuthenticationResponse(writer, 0x01, 0xFF) // nolint: errcheck
			return newError("invalid username or password")
		}

		if err := writeSocks5AuthenticationResponse(writer, 0x01, 0x00); err != nil {
			return newError("failed to write auth response").Base(err)
		}
	}

	return nil
}

func (s *ServerSession) handshake5(nMethod byte, reader io.Reader, writer io.Writer) (*protocol.RequestHeader, error) {
	if err := s.auth5(nMethod, reader, writer); err != nil {
		return nil, err
	}

	var cmd byte
	{
		buffer := buf.StackNew()
		if _, err := buffer.ReadFullFrom(reader, 3); err != nil {
			buffer.Release()
			return nil, newError("failed to read request").Base(err)
		}
		cmd = buffer.Byte(1)
		buffer.Release()
	}

	request := new(protocol.RequestHeader)
	switch cmd {
	case cmdTCPConnect, cmdTorResolve, cmdTorResolvePTR:
		// We don't have a solution for Tor case now. Simply treat it as connect command.
		request.Command = protocol.RequestCommandTCP
	case cmdUDPPort:
		if !s.config.UdpEnabled {
			writeSocks5Response(writer, statusCmdNotSupport, net.AnyIP, net.Port(0)) // nolint: errcheck
			return nil, newError("UDP is not enabled.")
		}
		request.Command = protocol.RequestCommandUDP
	case cmdTCPBind:
		writeSocks5Response(writer, statusCmdNotSupport, net.AnyIP, net.Port(0)) // nolint: errcheck
		return nil, newError("TCP bind is not supported.")
	default:
		writeSocks5Response(writer, statusCmdNotSupport, net.AnyIP, net.Port(0)) // nolint: errcheck
		return nil, newError("unknown command ", cmd)
	}

	request.Version = socks5Version

	addr, port, err := addrParser.ReadAddressPort(nil, reader)
	if err != nil {
		return nil, newError("failed to read address").Base(err)
	}
	request.Address = addr
	request.Port = port

	responseAddress := net.AnyIP
	responsePort := net.Port(1717)
	if request.Command == protocol.RequestCommandUDP {
		addr := s.config.Address.AsAddress()
		if addr == nil {
			addr = net.LocalHostIP
		}
		responseAddress = addr
		responsePort = s.port
	}
	if err := writeSocks5Response(writer, statusSuccess, responseAddress, responsePort); err != nil {
		return nil, err
	}

	return request, nil
}

// Handshake performs a Socks4/4a/5 handshake.
func (s *ServerSession) Handshake(reader io.Reader, writer io.Writer) (*protocol.RequestHeader, error) {
	buffer := buf.StackNew()
	if _, err := buffer.ReadFullFrom(reader, 2); err != nil {
		buffer.Release()
		return nil, newError("insufficient header").Base(err)
	}

	version := buffer.Byte(0)
	cmd := buffer.Byte(1)
	buffer.Release()

	switch version {
	case socks4Version:
		return s.handshake4(cmd, reader, writer)
	case socks5Version:
		return s.handshake5(cmd, reader, writer)
	default:
		return nil, newError("unknown Socks version: ", version)
	}
}

// ReadUsernamePassword reads Socks 5 username/password message from the given reader.
// +----+------+----------+------+----------+
// |VER | ULEN |  UNAME   | PLEN |  PASSWD  |
// +----+------+----------+------+----------+
// | 1  |  1   | 1 to 255 |  1   | 1 to 255 |
// +----+------+----------+------+----------+
func ReadUsernamePassword(reader io.Reader) (string, string, error) {
	buffer := buf.StackNew()
	defer buffer.Release()

	if _, err := buffer.ReadFullFrom(reader, 2); err != nil {
		return "", "", err
	}
	nUsername := int32(buffer.Byte(1))

	buffer.Clear()
	if _, err := buffer.ReadFullFrom(reader, nUsername); err != nil {
		return "", "", err
	}
	username := buffer.String()

	buffer.Clear()
	if _, err := buffer.ReadFullFrom(reader, 1); err != nil {
		return "", "", err
	}
	nPassword := int32(buffer.Byte(0))

	buffer.Clear()
	if _, err := buffer.ReadFullFrom(reader, nPassword); err != nil {
		return "", "", err
	}
	password := buffer.String()
	return username, password, nil
}

// ReadUntilNull reads content from given reader, until a null (0x00) byte.
func ReadUntilNull(reader io.Reader) (string, error) {
	b := buf.StackNew()
	defer b.Release()

	for {
		_, err := b.ReadFullFrom(reader, 1)
		if err != nil {
			return "", err
		}
		if b.Byte(b.Len()-1) == 0x00 {
			b.Resize(0, b.Len()-1)
			return b.String(), nil
		}
		if b.IsFull() {
			return "", newError("buffer overrun")
		}
	}
}

func hasAuthMethod(expectedAuth byte, authCandidates []byte) bool {
	for _, a := range authCandidates {
		if a == expectedAuth {
			return true
		}
	}
	return false
}

func writeSocks5AuthenticationResponse(writer io.Writer, version byte, auth byte) error {
	return buf.WriteAllBytes(writer, []byte{version, auth})
}

func writeSocks5Response(writer io.Writer, errCode byte, address net.Address, port net.Port) error {
	buffer := buf.New()
	defer buffer.Release()

	common.Must2(buffer.Write([]byte{socks5Version, errCode, 0x00 /* reserved */}))
	if err := addrParser.WriteAddressPort(buffer, address, port); err != nil {
		return err
	}

	return buf.WriteAllBytes(writer, buffer.Bytes())
}

func writeSocks4Response(writer io.Writer, errCode byte, address net.Address, port net.Port) error {
	buffer := buf.StackNew()
	defer buffer.Release()

	common.Must(buffer.WriteByte(0x00))
	common.Must(buffer.WriteByte(errCode))
	portBytes := buffer.Extend(2)
	binary.BigEndian.PutUint16(portBytes, port.Value())
	common.Must2(buffer.Write(address.IP()))
	return buf.WriteAllBytes(writer, buffer.Bytes())
}

func DecodeUDPPacket(packet *buf.Buffer) (*protocol.RequestHeader, error) {
	if packet.Len() < 5 {
		return nil, newError("insufficient length of packet.")
	}
	request := &protocol.RequestHeader{
		Version: socks5Version,
		Command: protocol.RequestCommandUDP,
	}

	// packet[0] and packet[1] are reserved
	if packet.Byte(2) != 0 /* fragments */ {
		return nil, newError("discarding fragmented payload.")
	}

	packet.Advance(3)

	addr, port, err := addrParser.ReadAddressPort(nil, packet)
	if err != nil {
		return nil, newError("failed to read UDP header").Base(err)
	}
	request.Address = addr
	request.Port = port
	return request, nil
}

func EncodeUDPPacket(request *protocol.RequestHeader, data []byte) (*buf.Buffer, error) {
	b := buf.New()
	common.Must2(b.Write([]byte{0, 0, 0 /* Fragment */}))
	if err := addrParser.WriteAddressPort(b, request.Address, request.Port); err != nil {
		b.Release()
		return nil, err
	}
	common.Must2(b.Write(data))
	return b, nil
}

type UDPReader struct {
	reader io.Reader
}

func NewUDPReader(reader io.Reader) *UDPReader {
	return &UDPReader{reader: reader}
}

func (r *UDPReader) ReadMultiBuffer() (buf.MultiBuffer, error) {
	b := buf.New()
	if _, err := b.ReadFrom(r.reader); err != nil {
		return nil, err
	}
	if _, err := DecodeUDPPacket(b); err != nil {
		return nil, err
	}
	return buf.MultiBuffer{b}, nil
}

type UDPWriter struct {
	request *protocol.RequestHeader
	writer  io.Writer
}

func NewUDPWriter(request *protocol.RequestHeader, writer io.Writer) *UDPWriter {
	return &UDPWriter{
		request: request,
		writer:  writer,
	}
}

// Write implements io.Writer.
func (w *UDPWriter) Write(b []byte) (int, error) {
	eb, err := EncodeUDPPacket(w.request, b)
	if err != nil {
		return 0, err
	}
	defer eb.Release()
	if _, err := w.writer.Write(eb.Bytes()); err != nil {
		return 0, err
	}
	return len(b), nil
}

func ClientHandshake(request *protocol.RequestHeader, reader io.Reader, writer io.Writer) (*protocol.RequestHeader, error) {
	authByte := byte(authNotRequired)
	if request.User != nil {
		authByte = byte(authPassword)
	}

	b := buf.New()
	defer b.Release()

	common.Must2(b.Write([]byte{socks5Version, 0x01, authByte}))
	if authByte == authPassword {
		account := request.User.Account.(*Account)

		common.Must(b.WriteByte(0x01))
		common.Must(b.WriteByte(byte(len(account.Username))))
		common.Must2(b.WriteString(account.Username))
		common.Must(b.WriteByte(byte(len(account.Password))))
		common.Must2(b.WriteString(account.Password))
	}

	if err := buf.WriteAllBytes(writer, b.Bytes()); err != nil {
		return nil, err
	}

	b.Clear()
	if _, err := b.ReadFullFrom(reader, 2); err != nil {
		return nil, err
	}

	if b.Byte(0) != socks5Version {
		return nil, newError("unexpected server version: ", b.Byte(0)).AtWarning()
	}
	if b.Byte(1) != authByte {
		return nil, newError("auth method not supported.").AtWarning()
	}

	if authByte == authPassword {
		b.Clear()
		if _, err := b.ReadFullFrom(reader, 2); err != nil {
			return nil, err
		}
		if b.Byte(1) != 0x00 {
			return nil, newError("server rejects account: ", b.Byte(1))
		}
	}

	b.Clear()

	command := byte(cmdTCPConnect)
	if request.Command == protocol.RequestCommandUDP {
		command = byte(cmdUDPPort)
	}
	common.Must2(b.Write([]byte{socks5Version, command, 0x00 /* reserved */}))
	if err := addrParser.WriteAddressPort(b, request.Address, request.Port); err != nil {
		return nil, err
	}

	if err := buf.WriteAllBytes(writer, b.Bytes()); err != nil {
		return nil, err
	}

	b.Clear()
	if _, err := b.ReadFullFrom(reader, 3); err != nil {
		return nil, err
	}

	resp := b.Byte(1)
	if resp != 0x00 {
		return nil, newError("server rejects request: ", resp)
	}

	b.Clear()

	address, port, err := addrParser.ReadAddressPort(b, reader)
	if err != nil {
		return nil, err
	}

	if request.Command == protocol.RequestCommandUDP {
		udpRequest := &protocol.RequestHeader{
			Version: socks5Version,
			Command: protocol.RequestCommandUDP,
			Address: address,
			Port:    port,
		}
		return udpRequest, nil
	}

	return nil, nil
}
