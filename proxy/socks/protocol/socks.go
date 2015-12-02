package protocol

import (
	"io"

	"github.com/v2ray/v2ray-core/common/alloc"
	"github.com/v2ray/v2ray-core/common/log"
	v2net "github.com/v2ray/v2ray-core/common/net"
	proxyerrors "github.com/v2ray/v2ray-core/proxy/common/errors"
	"github.com/v2ray/v2ray-core/transport"
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
	buffer := alloc.NewSmallBuffer()
	defer buffer.Release()

	nBytes, err := reader.Read(buffer.Value)
	if err != nil {
		return
	}
	if nBytes < 2 {
		log.Info("Socks expected 2 bytes read, but only %d bytes read", nBytes)
		err = transport.CorruptedPacket
		return
	}

	if buffer.Value[0] == socks4Version {
		auth4.Version = buffer.Value[0]
		auth4.Command = buffer.Value[1]
		auth4.Port = v2net.PortFromBytes(buffer.Value[2:4])
		copy(auth4.IP[:], buffer.Value[4:8])
		err = Socks4Downgrade
		return
	}

	auth.version = buffer.Value[0]
	if auth.version != socksVersion {
		log.Warning("Unknown protocol version %d", auth.version)
		err = proxyerrors.InvalidProtocolVersion
		return
	}

	auth.nMethods = buffer.Value[1]
	if auth.nMethods <= 0 {
		log.Info("Zero length of authentication methods")
		err = transport.CorruptedPacket
		return
	}

	if nBytes-2 != int(auth.nMethods) {
		log.Info("Unmatching number of auth methods, expecting %d, but got %d", auth.nMethods, nBytes)
		err = transport.CorruptedPacket
		return
	}
	copy(auth.authMethods[:], buffer.Value[2:nBytes])
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
	buffer := alloc.NewSmallBuffer()
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
	buffer := alloc.NewSmallBuffer()
	defer buffer.Release()

	nBytes, err := reader.Read(buffer.Value[:4])
	if err != nil {
		return
	}
	if nBytes < 4 {
		err = transport.CorruptedPacket
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
		nBytes, err = reader.Read(request.IPv4[:])
		if err != nil {
			return
		}
		if nBytes != 4 {
			err = transport.CorruptedPacket
			return
		}
	case AddrTypeDomain:
		nBytes, err = reader.Read(buffer.Value[0:1])
		if err != nil {
			return
		}
		domainLength := buffer.Value[0]
		nBytes, err = reader.Read(buffer.Value[:domainLength])
		if err != nil {
			return
		}

		if nBytes != int(domainLength) {
			log.Info("Unable to read domain with %d bytes, expecting %d bytes", nBytes, domainLength)
			err = transport.CorruptedPacket
			return
		}
		request.Domain = string(buffer.Value[:domainLength])
	case AddrTypeIPv6:
		nBytes, err = reader.Read(request.IPv6[:])
		if err != nil {
			return
		}
		if nBytes != 16 {
			err = transport.CorruptedPacket
			return
		}
	default:
		log.Info("Unexpected address type %d", request.AddrType)
		err = transport.CorruptedPacket
		return
	}

	nBytes, err = reader.Read(buffer.Value[:2])
	if err != nil {
		return
	}
	if nBytes != 2 {
		err = transport.CorruptedPacket
		return
	}

	request.Port = v2net.PortFromBytes(buffer.Value[:2])
	return
}

func (request *Socks5Request) Destination() v2net.Destination {
	var address v2net.Address
	switch request.AddrType {
	case AddrTypeIPv4:
		address = v2net.IPAddress(request.IPv4[:], request.Port)
	case AddrTypeIPv6:
		address = v2net.IPAddress(request.IPv6[:], request.Port)
	case AddrTypeDomain:
		address = v2net.DomainAddress(request.Domain, request.Port)
	default:
		panic("Unknown address type")
	}
	return v2net.NewTCPDestination(address)
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

func (r *Socks5Response) Write(buffer *alloc.Buffer) {
	buffer.AppendBytes(r.Version, r.Error, 0x00 /* reserved */, r.AddrType)
	switch r.AddrType {
	case 0x01:
		buffer.Append(r.IPv4[:])
	case 0x03:
		buffer.AppendBytes(byte(len(r.Domain)))
		buffer.Append([]byte(r.Domain))
	case 0x04:
		buffer.Append(r.IPv6[:])
	}
	buffer.Append(r.Port.Bytes())
}
