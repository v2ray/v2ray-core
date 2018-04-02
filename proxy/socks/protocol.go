package socks

import (
	"io"

	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/common/serial"
)

const (
	socks5Version = 0x05
	socks4Version = 0x04

	cmdTCPConnect = 0x01
	cmdTCPBind    = 0x02
	cmdUDPPort    = 0x03

	socks4RequestGranted  = 90
	socks4RequestRejected = 91

	authNotRequired = 0x00
	//authGssAPI           = 0x01
	authPassword         = 0x02
	authNoMatchingMethod = 0xFF

	addrTypeIPv4   = 0x01
	addrTypeIPv6   = 0x04
	addrTypeDomain = 0x03

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

func (s *ServerSession) Handshake(reader io.Reader, writer io.Writer) (*protocol.RequestHeader, error) {
	buffer := buf.New()
	defer buffer.Release()

	request := new(protocol.RequestHeader)

	if err := buffer.AppendSupplier(buf.ReadFullFrom(reader, 2)); err != nil {
		return nil, newError("insufficient header").Base(err)
	}

	version := buffer.Byte(0)
	if version == socks4Version {
		if s.config.AuthType == AuthType_PASSWORD {
			writeSocks4Response(writer, socks4RequestRejected, net.AnyIP, net.Port(0))
			return nil, newError("socks 4 is not allowed when auth is required.")
		}

		if err := buffer.AppendSupplier(buf.ReadFullFrom(reader, 6)); err != nil {
			return nil, newError("insufficient header").Base(err)
		}
		port := net.PortFromBytes(buffer.BytesRange(2, 4))
		address := net.IPAddress(buffer.BytesRange(4, 8))
		_, err := readUntilNull(reader) // user id
		if err != nil {
			return nil, err
		}
		if address.IP()[0] == 0x00 {
			domain, err := readUntilNull(reader)
			if err != nil {
				return nil, newError("failed to read domain for socks 4a").Base(err)
			}
			address = net.DomainAddress(domain)
		}

		switch buffer.Byte(1) {
		case cmdTCPConnect:
			request.Command = protocol.RequestCommandTCP
			request.Address = address
			request.Port = port
			request.Version = socks4Version
			if err := writeSocks4Response(writer, socks4RequestGranted, net.AnyIP, net.Port(0)); err != nil {
				return nil, err
			}
			return request, nil
		default:
			writeSocks4Response(writer, socks4RequestRejected, net.AnyIP, net.Port(0))
			return nil, newError("unsupported command: ", buffer.Byte(1))
		}
	}

	if version == socks5Version {
		nMethod := int(buffer.Byte(1))
		if err := buffer.AppendSupplier(buf.ReadFullFrom(reader, int32(nMethod))); err != nil {
			return nil, newError("failed to read auth methods").Base(err)
		}

		var expectedAuth byte = authNotRequired
		if s.config.AuthType == AuthType_PASSWORD {
			expectedAuth = authPassword
		}

		if !hasAuthMethod(expectedAuth, buffer.BytesRange(2, 2+nMethod)) {
			writeSocks5AuthenticationResponse(writer, socks5Version, authNoMatchingMethod)
			return nil, newError("no matching auth method")
		}

		if err := writeSocks5AuthenticationResponse(writer, socks5Version, expectedAuth); err != nil {
			return nil, newError("failed to write auth response").Base(err)
		}

		if expectedAuth == authPassword {
			username, password, err := readUsernamePassword(reader)
			if err != nil {
				return nil, newError("failed to read username and password for authentication").Base(err)
			}

			if !s.config.HasAccount(username, password) {
				writeSocks5AuthenticationResponse(writer, 0x01, 0xFF)
				return nil, newError("invalid username or password")
			}

			if err := writeSocks5AuthenticationResponse(writer, 0x01, 0x00); err != nil {
				return nil, newError("failed to write auth response").Base(err)
			}
		}
		if err := buffer.Reset(buf.ReadFullFrom(reader, 3)); err != nil {
			return nil, newError("failed to read request").Base(err)
		}

		cmd := buffer.Byte(1)
		switch cmd {
		case cmdTCPConnect:
			request.Command = protocol.RequestCommandTCP
		case cmdUDPPort:
			if !s.config.UdpEnabled {
				writeSocks5Response(writer, statusCmdNotSupport, net.AnyIP, net.Port(0))
				return nil, newError("UDP is not enabled.")
			}
			request.Command = protocol.RequestCommandUDP
		case cmdTCPBind:
			writeSocks5Response(writer, statusCmdNotSupport, net.AnyIP, net.Port(0))
			return nil, newError("TCP bind is not supported.")
		default:
			writeSocks5Response(writer, statusCmdNotSupport, net.AnyIP, net.Port(0))
			return nil, newError("unknown command ", cmd)
		}

		buffer.Clear()

		request.Version = socks5Version

		addr, port, err := addrParser.ReadAddressPort(buffer, reader)
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

	return nil, newError("unknown Socks version: ", version)
}

func readUsernamePassword(reader io.Reader) (string, string, error) {
	buffer := buf.New()
	defer buffer.Release()

	if err := buffer.Reset(buf.ReadFullFrom(reader, 2)); err != nil {
		return "", "", err
	}
	nUsername := int32(buffer.Byte(1))

	if err := buffer.Reset(buf.ReadFullFrom(reader, nUsername)); err != nil {
		return "", "", err
	}
	username := buffer.String()

	if err := buffer.Reset(buf.ReadFullFrom(reader, 1)); err != nil {
		return "", "", err
	}
	nPassword := int32(buffer.Byte(0))
	if err := buffer.Reset(buf.ReadFullFrom(reader, nPassword)); err != nil {
		return "", "", err
	}
	password := buffer.String()
	return username, password, nil
}

func readUntilNull(reader io.Reader) (string, error) {
	var b [256]byte
	size := 0
	for {
		_, err := reader.Read(b[size : size+1])
		if err != nil {
			return "", err
		}
		if b[size] == 0x00 {
			return string(b[:size]), nil
		}
		size++
		if size == 256 {
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
	_, err := writer.Write([]byte{version, auth})
	return err
}

func writeSocks5Response(writer io.Writer, errCode byte, address net.Address, port net.Port) error {
	buffer := buf.New()
	defer buffer.Release()

	buffer.AppendBytes(socks5Version, errCode, 0x00 /* reserved */)
	if err := addrParser.WriteAddressPort(buffer, address, port); err != nil {
		return err
	}

	_, err := writer.Write(buffer.Bytes())
	return err
}

func writeSocks4Response(writer io.Writer, errCode byte, address net.Address, port net.Port) error {
	buffer := buf.New()
	defer buffer.Release()

	buffer.AppendBytes(0x00, errCode)
	common.Must(buffer.AppendSupplier(serial.WriteUint16(port.Value())))
	buffer.Append(address.IP())
	_, err := writer.Write(buffer.Bytes())
	return err
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

	packet.SliceFrom(3)

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
	b.AppendBytes(0, 0, 0 /* Fragment */)
	if err := addrParser.WriteAddressPort(b, request.Address, request.Port); err != nil {
		b.Release()
		return nil, err
	}
	b.Append(data)
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
	if err := b.AppendSupplier(buf.ReadFrom(r.reader)); err != nil {
		return nil, err
	}
	if _, err := DecodeUDPPacket(b); err != nil {
		return nil, err
	}
	return buf.NewMultiBufferValue(b), nil
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

	b.AppendBytes(socks5Version, 0x01, authByte)
	if authByte == authPassword {
		rawAccount, err := request.User.GetTypedAccount()
		if err != nil {
			return nil, err
		}
		account := rawAccount.(*Account)

		b.AppendBytes(0x01, byte(len(account.Username)))
		b.Append([]byte(account.Username))
		b.AppendBytes(byte(len(account.Password)))
		b.Append([]byte(account.Password))
	}

	if _, err := writer.Write(b.Bytes()); err != nil {
		return nil, err
	}

	if err := b.Reset(buf.ReadFullFrom(reader, 2)); err != nil {
		return nil, err
	}

	if b.Byte(0) != socks5Version {
		return nil, newError("unexpected server version: ", b.Byte(0)).AtWarning()
	}
	if b.Byte(1) != authByte {
		return nil, newError("auth method not supported.").AtWarning()
	}

	if authByte == authPassword {
		if err := b.Reset(buf.ReadFullFrom(reader, 2)); err != nil {
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
	b.AppendBytes(socks5Version, command, 0x00 /* reserved */)
	if err := addrParser.WriteAddressPort(b, request.Address, request.Port); err != nil {
		return nil, err
	}

	if _, err := writer.Write(b.Bytes()); err != nil {
		return nil, err
	}

	b.Clear()
	if err := b.AppendSupplier(buf.ReadFullFrom(reader, 3)); err != nil {
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
