package protocol

import (
	"io"

	"v2ray.com/core/common/alloc"
	"v2ray.com/core/common/log"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/proxy"
	"v2ray.com/core/transport"
)

const (
	socksVersion  = byte(0x05)
	socks4Version = byte(0x04)

	AuthNotRequired      = byte(0x00)
	AuthGssApi           = byte(0x01)
	AuthUserPass         = byte(0x02)
	AuthNoMatchingMethod = byte(0xFF)

	Socks4RequestGranted  = byte(90)
	Socks4RequestRejected = byte(91)
)

// Authentication request header of Socks5 protocol
type Socks5AuthenticationRequest struct {
	version     byte
	nMethods    byte
	authMethods [256]byte
}

func (request *Socks5AuthenticationRequest) HasAuthMethod(method byte) bool {
	for i := 0; i < int(request.nMethods); i++ {
		if request.authMethods[i] == method {
			return true
		}
	}
	return false
}

func ReadAuthentication(reader io.Reader) (auth Socks5AuthenticationRequest, auth4 Socks4AuthenticationRequest, err error) {
	buffer := make([]byte, 256)

	nBytes, err := reader.Read(buffer)
	if err != nil {
		return
	}
	if nBytes < 2 {
		log.Warning("Socks: expected 2 bytes read, but only ", nBytes, " bytes read")
		err = transport.ErrCorruptedPacket
		return
	}

	if buffer[0] == socks4Version {
		auth4.Version = buffer[0]
		auth4.Command = buffer[1]
		auth4.Port = v2net.PortFromBytes(buffer[2:4])
		copy(auth4.IP[:], buffer[4:8])
		err = Socks4Downgrade
		return
	}

	auth.version = buffer[0]
	if auth.version != socksVersion {
		log.Warning("Socks: Unknown protocol version ", auth.version)
		err = proxy.ErrInvalidProtocolVersion
		return
	}

	auth.nMethods = buffer[1]
	if auth.nMethods <= 0 {
		log.Warning("Socks: Zero length of authentication methods")
		err = proxy.ErrInvalidAuthentication
		return
	}

	if nBytes-2 != int(auth.nMethods) {
		log.Warning("Socks: Unmatching number of auth methods, expecting ", auth.nMethods, ", but got ", nBytes)
		err = proxy.ErrInvalidAuthentication
		return
	}
	copy(auth.authMethods[:], buffer[2:nBytes])
	return
}

type Socks5AuthenticationResponse struct {
	version    byte
	authMethod byte
}

func NewAuthenticationResponse(authMethod byte) *Socks5AuthenticationResponse {
	return &Socks5AuthenticationResponse{
		version:    socksVersion,
		authMethod: authMethod,
	}
}

func WriteAuthentication(writer io.Writer, r *Socks5AuthenticationResponse) error {
	_, err := writer.Write([]byte{r.version, r.authMethod})
	return err
}

type Socks5UserPassRequest struct {
	version  byte
	username string
	password string
}

func (request Socks5UserPassRequest) Username() string {
	return request.username
}

func (request Socks5UserPassRequest) Password() string {
	return request.password
}

func (request Socks5UserPassRequest) AuthDetail() string {
	return request.username + ":" + request.password
}

func ReadUserPassRequest(reader io.Reader) (request Socks5UserPassRequest, err error) {
	buffer := alloc.NewLocalBuffer(512)
	defer buffer.Release()

	_, err = reader.Read(buffer.Value[0:2])
	if err != nil {
		return
	}
	request.version = buffer.Value[0]
	nUsername := buffer.Value[1]
	nBytes, err := reader.Read(buffer.Value[:nUsername])
	if err != nil {
		return
	}
	request.username = string(buffer.Value[:nBytes])

	_, err = reader.Read(buffer.Value[0:1])
	if err != nil {
		return
	}
	nPassword := buffer.Value[0]
	nBytes, err = reader.Read(buffer.Value[:nPassword])
	if err != nil {
		return
	}
	request.password = string(buffer.Value[:nBytes])
	return
}

type Socks5UserPassResponse struct {
	version byte
	status  byte
}

func NewSocks5UserPassResponse(status byte) Socks5UserPassResponse {
	return Socks5UserPassResponse{
		version: socksVersion,
		status:  status,
	}
}

func WriteUserPassResponse(writer io.Writer, response Socks5UserPassResponse) error {
	_, err := writer.Write([]byte{response.version, response.status})
	return err
}

const (
	AddrTypeIPv4   = byte(0x01)
	AddrTypeIPv6   = byte(0x04)
	AddrTypeDomain = byte(0x03)

	CmdConnect      = byte(0x01)
	CmdBind         = byte(0x02)
	CmdUdpAssociate = byte(0x03)
)

type Socks5Request struct {
	Version  byte
	Command  byte
	AddrType byte
	IPv4     [4]byte
	Domain   string
	IPv6     [16]byte
	Port     v2net.Port
}

func ReadRequest(reader io.Reader) (request *Socks5Request, err error) {
	buffer := alloc.NewLocalBuffer(512)
	defer buffer.Release()

	_, err = io.ReadFull(reader, buffer.Value[:4])
	if err != nil {
		return
	}

	request = &Socks5Request{
		Version: buffer.Value[0],
		Command: buffer.Value[1],
		// buffer[2] is a reserved field
		AddrType: buffer.Value[3],
	}
	switch request.AddrType {
	case AddrTypeIPv4:
		_, err = io.ReadFull(reader, request.IPv4[:])
		if err != nil {
			return
		}
	case AddrTypeDomain:
		_, err = io.ReadFull(reader, buffer.Value[0:1])
		if err != nil {
			return
		}
		domainLength := buffer.Value[0]
		_, err = io.ReadFull(reader, buffer.Value[:domainLength])
		if err != nil {
			return
		}

		request.Domain = string(append([]byte(nil), buffer.Value[:domainLength]...))
	case AddrTypeIPv6:
		_, err = io.ReadFull(reader, request.IPv6[:])
		if err != nil {
			return
		}
	default:
		log.Warning("Socks: Unexpected address type ", request.AddrType)
		err = transport.ErrCorruptedPacket
		return
	}

	_, err = io.ReadFull(reader, buffer.Value[:2])
	if err != nil {
		return
	}

	request.Port = v2net.PortFromBytes(buffer.Value[:2])
	return
}

func (request *Socks5Request) Destination() v2net.Destination {
	switch request.AddrType {
	case AddrTypeIPv4:
		return v2net.TCPDestination(v2net.IPAddress(request.IPv4[:]), request.Port)
	case AddrTypeIPv6:
		return v2net.TCPDestination(v2net.IPAddress(request.IPv6[:]), request.Port)
	case AddrTypeDomain:
		return v2net.TCPDestination(v2net.ParseAddress(request.Domain), request.Port)
	default:
		panic("Unknown address type")
	}
}

const (
	ErrorSuccess                 = byte(0x00)
	ErrorGeneralFailure          = byte(0x01)
	ErrorConnectionNotAllowed    = byte(0x02)
	ErrorNetworkUnreachable      = byte(0x03)
	ErrorHostUnUnreachable       = byte(0x04)
	ErrorConnectionRefused       = byte(0x05)
	ErrorTTLExpired              = byte(0x06)
	ErrorCommandNotSupported     = byte(0x07)
	ErrorAddressTypeNotSupported = byte(0x08)
)

type Socks5Response struct {
	Version  byte
	Error    byte
	AddrType byte
	IPv4     [4]byte
	Domain   string
	IPv6     [16]byte
	Port     v2net.Port
}

func NewSocks5Response() *Socks5Response {
	return &Socks5Response{
		Version: socksVersion,
	}
}

func (r *Socks5Response) SetIPv4(ipv4 []byte) {
	r.AddrType = AddrTypeIPv4
	copy(r.IPv4[:], ipv4)
}

func (r *Socks5Response) SetIPv6(ipv6 []byte) {
	r.AddrType = AddrTypeIPv6
	copy(r.IPv6[:], ipv6)
}

func (r *Socks5Response) SetDomain(domain string) {
	r.AddrType = AddrTypeDomain
	r.Domain = domain
}

func (r *Socks5Response) Write(writer io.Writer) {
	writer.Write([]byte{r.Version, r.Error, 0x00 /* reserved */, r.AddrType})
	switch r.AddrType {
	case 0x01:
		writer.Write(r.IPv4[:])
	case 0x03:
		writer.Write([]byte{byte(len(r.Domain))})
		writer.Write([]byte(r.Domain))
	case 0x04:
		writer.Write(r.IPv6[:])
	}
	writer.Write(r.Port.Bytes(nil))
}
