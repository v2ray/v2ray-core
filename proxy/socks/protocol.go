package socks

import (
	"io"

	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/errors"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/common/serial"
	"v2ray.com/core/proxy"
)

const (
	socks5Version = 0x05
	socks4Version = 0x04

	cmdTCPConnect = 0x01
	cmdTCPBind    = 0x02
	cmdUDPPort    = 0x03

	socks4RequestGranted  = 90
	socks4RequestRejected = 91

	authNotRequired      = 0x00
	authGssAPI           = 0x01
	authPassword         = 0x02
	authNoMatchingMethod = 0xFF

	addrTypeIPv4   = 0x01
	addrTypeIPv6   = 0x04
	addrTypeDomain = 0x03

	statusSuccess       = 0x00
	statusCmdNotSupport = 0x07
)

type ServerSession struct {
	config *ServerConfig
	meta   *proxy.InboundHandlerMeta
}

func (s *ServerSession) Handshake(reader io.Reader, writer io.Writer) (*protocol.RequestHeader, error) {
	buffer := buf.NewLocal(512)
	request := new(protocol.RequestHeader)

	if err := buffer.AppendSupplier(buf.ReadFullFrom(reader, 2)); err != nil {
		return nil, errors.Base(err).Message("Socks|Server: Insufficient header.")
	}

	version := buffer.Byte(0)
	if version == socks4Version {
		if err := buffer.AppendSupplier(buf.ReadFullFrom(reader, 6)); err != nil {
			return nil, errors.Base(err).Message("Socks|Server: Insufficient header.")
		}
		port := v2net.PortFromBytes(buffer.BytesRange(2, 4))
		address := v2net.IPAddress(buffer.BytesRange(4, 8))
		_, err := readUntilNull(reader) // user id
		if err != nil {
			return nil, err
		}
		if address.IP()[0] == 0x00 {
			domain, err := readUntilNull(reader)
			if err != nil {
				return nil, errors.Base(err).Message("Socks|Server: Failed to read domain for socks 4a.")
			}
			address = v2net.DomainAddress(domain)
		}

		switch buffer.Byte(1) {
		case cmdTCPConnect:
			request.Command = protocol.RequestCommandTCP
			request.Address = address
			request.Port = port
			request.Version = socks4Version
			if err := writeSocks4Response(writer, socks4RequestGranted, address, port); err != nil {
				return nil, err
			}
			return request, nil
		default:
			writeSocks4Response(writer, socks4RequestRejected, address, port)
			return nil, errors.New("Socks|Server: Unsupported command: ", buffer.Byte(1))
		}
	}

	if version == socks5Version {
		nMethod := int(buffer.Byte(1))
		if err := buffer.AppendSupplier(buf.ReadFullFrom(reader, nMethod)); err != nil {
			return nil, err
		}

		var expectedAuth byte = authNotRequired
		if len(s.config.Accounts) > 0 {
			expectedAuth = authPassword
		}

		if !hasAuthMethod(expectedAuth, buffer.BytesRange(2, 2+nMethod)) {
			writeSocks5AuthenticationResponse(writer, authNoMatchingMethod)
			return nil, errors.New("Socks|Server: No matching auth method.")
		}

		if expectedAuth == authPassword {
			username, password, err := readUsernamePassword(reader)
			if err != nil {
				return nil, errors.Base(err).Message("Socks|Server: Failed to read username or password.")
			}
			if !s.config.HasAccount(username, password) {
				writeSocks5AuthenticationResponse(writer, 0xFF)
				return nil, errors.Base(err).Message("Socks|Server: Invalid username or password.")
			}
		}

		if err := writeSocks5AuthenticationResponse(writer, 0x00); err != nil {
			return nil, err
		}

		buffer.Clear()
		if err := buffer.AppendSupplier(buf.ReadFullFrom(reader, 4)); err != nil {
			return nil, err
		}

		cmd := buffer.Byte(1)
		if cmd == cmdTCPBind || (cmd == cmdUDPPort && !s.config.UdpEnabled) {
			writeSocks5Response(writer, statusCmdNotSupport, v2net.AnyIP, v2net.Port(0))
			return nil, errors.New("Socks|Server: Unsupported command: ", cmd)
		}

		switch cmd {
		case cmdTCPConnect:
			request.Command = protocol.RequestCommandTCP
		case cmdUDPPort:
			request.Command = protocol.RequestCommandUDP
		}

		addrType := buffer.Byte(3)

		buffer.Clear()

		request.Version = socks5Version
		switch addrType {
		case addrTypeIPv4:
			if err := buffer.AppendSupplier(buf.ReadFullFrom(reader, 4)); err != nil {
				return nil, err
			}
			request.Address = v2net.IPAddress(buffer.Bytes())
		case addrTypeIPv6:
			if err := buffer.AppendSupplier(buf.ReadFullFrom(reader, 16)); err != nil {
				return nil, err
			}
			request.Address = v2net.IPAddress(buffer.Bytes())
		case addrTypeDomain:
			if err := buffer.AppendSupplier(buf.ReadFullFrom(reader, 1)); err != nil {
				return nil, err
			}
			domainLength := int(buffer.Byte(0))
			if err := buffer.AppendSupplier(buf.ReadFullFrom(reader, domainLength)); err != nil {
				return nil, err
			}
			request.Address = v2net.DomainAddress(string(buffer.BytesFrom(-domainLength)))
		default:
			return nil, errors.New("Socks|Server: Unknown address type: ", addrType)
		}

		if err := buffer.AppendSupplier(buf.ReadFullFrom(reader, 2)); err != nil {
			return nil, err
		}
		request.Port = v2net.PortFromBytes(buffer.BytesFrom(-2))

		responseAddress := v2net.AnyIP
		responsePort := v2net.Port(1717)
		if request.Command == protocol.RequestCommandUDP {
			addr := s.config.Address.AsAddress()
			if addr == nil {
				addr = v2net.LocalHostIP
			}
			responseAddress = addr
			responsePort = s.meta.Port
		}
		if err := writeSocks5Response(writer, statusSuccess, responseAddress, responsePort); err != nil {
			return nil, err
		}

		return request, nil
	}

	return nil, errors.New("Socks|Server: Unknown Socks version: ", version)
}

func readUsernamePassword(reader io.Reader) (string, string, error) {
	buffer := buf.NewLocal(512)
	defer buffer.Release()

	if err := buffer.AppendSupplier(buf.ReadFullFrom(reader, 2)); err != nil {
		return "", "", err
	}
	nUsername := int(buffer.Byte(1))

	buffer.Clear()
	if err := buffer.AppendSupplier(buf.ReadFullFrom(reader, nUsername)); err != nil {
		return "", "", err
	}
	username := buffer.String()

	if err := buffer.AppendSupplier(buf.ReadFullFrom(reader, 1)); err != nil {
		return "", "", err
	}
	nPassword := int(buffer.Byte(0))
	if err := buffer.AppendSupplier(buf.ReadFullFrom(reader, nPassword)); err != nil {
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
			return "", errors.New("Socks|Server: Buffer overrun.")
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

func writeSocks5AuthenticationResponse(writer io.Writer, auth byte) error {
	_, err := writer.Write([]byte{socks5Version, auth})
	return err
}

func writeSocks5Response(writer io.Writer, errCode byte, address v2net.Address, port v2net.Port) error {
	buffer := buf.NewLocal(64)
	buffer.AppendBytes(socks5Version, errCode, 0x00 /* reserved */)
	switch address.Family() {
	case v2net.AddressFamilyIPv4:
		buffer.AppendBytes(0x01)
		buffer.Append(address.IP())
	case v2net.AddressFamilyIPv6:
		buffer.AppendBytes(0x04)
		buffer.Append(address.IP())
	case v2net.AddressFamilyDomain:
		buffer.AppendBytes(0x03, byte(len(address.Domain())))
		buffer.AppendSupplier(serial.WriteString(address.Domain()))
	}
	buffer.AppendSupplier(serial.WriteUint16(port.Value()))

	_, err := writer.Write(buffer.Bytes())
	return err
}

func writeSocks4Response(writer io.Writer, errCode byte, address v2net.Address, port v2net.Port) error {
	buffer := buf.NewLocal(32)
	buffer.AppendBytes(0x00, errCode)
	buffer.AppendSupplier(serial.WriteUint16(port.Value()))
	buffer.Append(address.IP())
	_, err := writer.Write(buffer.Bytes())
	return err
}

func DecodeUDPPacket(packet []byte) (*protocol.RequestHeader, []byte, error) {
	if len(packet) < 5 {
		return nil, nil, errors.New("Socks|UDP: Insufficient length of packet.")
	}
	request := &protocol.RequestHeader{
		Version: socks5Version,
		Command: protocol.RequestCommandUDP,
	}

	// packet[0] and packet[1] are reserved
	if packet[2] != 0 /* fragments */ {
		return nil, nil, errors.New("Socks|UDP: Fragmented payload.")
	}

	addrType := packet[3]
	var dataBegin int

	switch addrType {
	case addrTypeIPv4:
		if len(packet) < 10 {
			return nil, nil, errors.New("Socks|UDP: Insufficient length of packet.")
		}
		ip := packet[4:8]
		request.Port = v2net.PortFromBytes(packet[8:10])
		request.Address = v2net.IPAddress(ip)
		dataBegin = 10
	case addrTypeIPv6:
		if len(packet) < 22 {
			return nil, nil, errors.New("Socks|UDP: Insufficient length of packet.")
		}
		ip := packet[4:20]
		request.Port = v2net.PortFromBytes(packet[20:22])
		request.Address = v2net.IPAddress(ip)
		dataBegin = 22
	case addrTypeDomain:
		domainLength := int(packet[4])
		if len(packet) < 5+domainLength+2 {
			return nil, nil, errors.New("Socks|UDP: Insufficient length of packet.")
		}
		domain := string(packet[5 : 5+domainLength])
		request.Port = v2net.PortFromBytes(packet[5+domainLength : 5+domainLength+2])
		request.Address = v2net.ParseAddress(domain)
		dataBegin = 5 + domainLength + 2
	default:
		return nil, nil, errors.New("Socks|UDP: Unknown address type ", addrType)
	}

	return request, packet[dataBegin:], nil
}

func EncodeUDPPacket(request *protocol.RequestHeader, data []byte) *buf.Buffer {
	b := buf.NewSmall()
	b.AppendBytes(0, 0, 0 /* Fragment */)
	switch request.Address.Family() {
	case v2net.AddressFamilyIPv4:
		b.AppendBytes(addrTypeIPv4)
		b.Append(request.Address.IP())
	case v2net.AddressFamilyIPv6:
		b.AppendBytes(addrTypeIPv6)
		b.Append(request.Address.IP())
	case v2net.AddressFamilyDomain:
		b.AppendBytes(addrTypeDomain, byte(len(request.Address.Domain())))
		b.AppendSupplier(serial.WriteString(request.Address.Domain()))
	}
	b.AppendSupplier(serial.WriteUint16(request.Port.Value()))
	b.Append(data)
	return b
}
