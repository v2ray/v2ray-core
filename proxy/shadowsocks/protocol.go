package shadowsocks

import (
	"bytes"
	"crypto/rand"
	"io"
	"v2ray.com/core/common/alloc"
	"v2ray.com/core/common/crypto"
	"v2ray.com/core/common/errors"
	v2io "v2ray.com/core/common/io"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/common/protocol"
)

const (
	Version                  = 1
	RequestOptionOneTimeAuth = protocol.RequestOption(101)

	AddrTypeIPv4   = 1
	AddrTypeIPv6   = 4
	AddrTypeDomain = 3
)

func ReadTCPSession(user *protocol.User, reader io.Reader) (*protocol.RequestHeader, v2io.Reader, error) {
	rawAccount, err := user.GetTypedAccount()
	if err != nil {
		return nil, nil, errors.Base(err).Message("Shadowsocks|TCP: Failed to parse account.")
	}
	account := rawAccount.(*ShadowsocksAccount)

	buffer := alloc.NewLocalBuffer(512).Clear()
	defer buffer.Release()

	ivLen := account.Cipher.IVSize()
	_, err = buffer.FillFullFrom(reader, ivLen)
	if err != nil {
		return nil, nil, errors.Base(err).Message("Shadowsocks|TCP: Failed to read IV.")
	}

	iv := append([]byte(nil), buffer.Bytes()[:ivLen]...)

	stream, err := account.Cipher.NewDecodingStream(account.Key, iv)
	if err != nil {
		return nil, nil, errors.Base(err).Message("Shadowsocks|TCP: Failed to initialize decoding stream.")
	}
	reader = crypto.NewCryptionReader(stream, reader)

	authenticator := NewAuthenticator(HeaderKeyGenerator(account.Key, iv))
	request := &protocol.RequestHeader{
		Version: Version,
		User:    user,
		Command: protocol.RequestCommandTCP,
	}

	buffer.Clear()
	_, err = buffer.FillFullFrom(reader, 1)
	if err != nil {
		return nil, nil, errors.Base(err).Message("Shadowsocks|TCP: Failed to read address type.")
	}

	addrType := (buffer.Bytes()[0] & 0x0F)
	if (buffer.Bytes()[0] & 0x10) == 0x10 {
		request.Option |= RequestOptionOneTimeAuth
	}

	if request.Option.Has(RequestOptionOneTimeAuth) && account.OneTimeAuth == Account_Disabled {
		return nil, nil, errors.New("Shadowsocks|TCP: Rejecting connection with OTA enabled, while server disables OTA.")
	}

	if !request.Option.Has(RequestOptionOneTimeAuth) && account.OneTimeAuth == Account_Enabled {
		return nil, nil, errors.New("Shadowsocks|TCP: Rejecting connection with OTA disabled, while server enables OTA.")
	}

	switch addrType {
	case AddrTypeIPv4:
		_, err := buffer.FillFullFrom(reader, 4)
		if err != nil {
			return nil, nil, errors.Base(err).Message("Shadowsocks|TCP: Failed to read IPv4 address.")
		}
		request.Address = v2net.IPAddress(buffer.BytesFrom(-4))
	case AddrTypeIPv6:
		_, err := buffer.FillFullFrom(reader, 16)
		if err != nil {
			return nil, nil, errors.Base(err).Message("Shadowsocks|TCP: Failed to read IPv6 address.")
		}
		request.Address = v2net.IPAddress(buffer.BytesFrom(-16))
	case AddrTypeDomain:
		_, err := buffer.FillFullFrom(reader, 1)
		if err != nil {
			return nil, nil, errors.Base(err).Message("Shadowsocks|TCP: Failed to read domain lenth.")
		}
		domainLength := int(buffer.BytesFrom(-1)[0])
		_, err = buffer.FillFullFrom(reader, domainLength)
		if err != nil {
			return nil, nil, errors.Base(err).Message("Shadowsocks|TCP: Failed to read domain.")
		}
		request.Address = v2net.DomainAddress(string(buffer.BytesFrom(-domainLength)))
	default:
		return nil, nil, errors.New("Shadowsocks|TCP: Unknown address type: ", addrType)
	}

	_, err = buffer.FillFullFrom(reader, 2)
	if err != nil {
		return nil, nil, errors.Base(err).Message("Shadowsocks|TCP: Failed to read port.")
	}
	request.Port = v2net.PortFromBytes(buffer.BytesFrom(-2))

	if request.Option.Has(RequestOptionOneTimeAuth) {
		actualAuth := authenticator.Authenticate(nil, buffer.Bytes())

		_, err := buffer.FillFullFrom(reader, AuthSize)
		if err != nil {
			return nil, nil, errors.Base(err).Message("Shadowsocks|TCP: Failed to read OTA.")
		}

		if !bytes.Equal(actualAuth, buffer.BytesFrom(-AuthSize)) {
			return nil, nil, errors.New("Shadowsocks|TCP: Invalid OTA")
		}
	}

	var chunkReader v2io.Reader
	if request.Option.Has(RequestOptionOneTimeAuth) {
		chunkReader = NewChunkReader(reader, NewAuthenticator(ChunkKeyGenerator(iv)))
	} else {
		chunkReader = v2io.NewAdaptiveReader(reader)
	}

	return request, chunkReader, nil
}

func WriteTCPRequest(request *protocol.RequestHeader, writer io.Writer) (v2io.Writer, error) {
	user := request.User
	rawAccount, err := user.GetTypedAccount()
	if err != nil {
		return nil, errors.Base(err).Message("Shadowsocks|TCP: Failed to parse account.")
	}
	account := rawAccount.(*ShadowsocksAccount)

	iv := make([]byte, account.Cipher.IVSize())
	rand.Read(iv)
	_, err = writer.Write(iv)
	if err != nil {
		return nil, errors.Base(err).Message("Shadowsocks|TCP: Failed to write IV.")
	}

	stream, err := account.Cipher.NewEncodingStream(account.Key, iv)
	if err != nil {
		return nil, errors.Base(err).Message("Shadowsocks|TCP: Failed to create encoding stream.")
	}

	writer = crypto.NewCryptionWriter(stream, writer)

	header := alloc.NewLocalBuffer(512).Clear()

	switch request.Address.Family() {
	case v2net.AddressFamilyIPv4:
		header.AppendBytes(AddrTypeIPv4)
		header.Append([]byte(request.Address.IP()))
	case v2net.AddressFamilyIPv6:
		header.AppendBytes(AddrTypeIPv6)
		header.Append([]byte(request.Address.IP()))
	case v2net.AddressFamilyDomain:
		header.AppendBytes(AddrTypeDomain, byte(len(request.Address.Domain())))
		header.Append([]byte(request.Address.Domain()))
	default:
		return nil, errors.New("Shadowsocks|TCP: Unsupported address type: ", request.Address.Family())
	}

	header.AppendUint16(uint16(request.Port))

	if request.Option.Has(RequestOptionOneTimeAuth) {
		header.Bytes()[0] |= 0x10

		authenticator := NewAuthenticator(HeaderKeyGenerator(account.Key, iv))
		header.Value = authenticator.Authenticate(header.Value, header.Value)
	}

	_, err = writer.Write(header.Value)
	if err != nil {
		return nil, errors.Base(err).Message("Shadowsocks|TCP: Failed to write header.")
	}

	var chunkWriter v2io.Writer
	if request.Option.Has(RequestOptionOneTimeAuth) {
		chunkWriter = NewChunkWriter(writer, NewAuthenticator(ChunkKeyGenerator(iv)))
	} else {
		chunkWriter = v2io.NewAdaptiveWriter(writer)
	}

	return chunkWriter, nil
}

func ReadTCPResponse(user *protocol.User, reader io.Reader) (v2io.Reader, error) {
	rawAccount, err := user.GetTypedAccount()
	if err != nil {
		return nil, errors.Base(err).Message("Shadowsocks|TCP: Failed to parse account.")
	}
	account := rawAccount.(*ShadowsocksAccount)

	iv := make([]byte, account.Cipher.IVSize())
	_, err = io.ReadFull(reader, iv)
	if err != nil {
		return nil, errors.Base(err).Message("Shadowsocks|TCP: Failed to read IV.")
	}

	stream, err := account.Cipher.NewDecodingStream(account.Key, iv)
	if err != nil {
		return nil, errors.Base(err).Message("Shadowsocks|TCP: Failed to initialize decoding stream.")
	}
	return v2io.NewAdaptiveReader(crypto.NewCryptionReader(stream, reader)), nil
}

func WriteTCPResponse(request *protocol.RequestHeader, writer io.Writer) (v2io.Writer, error) {
	user := request.User
	rawAccount, err := user.GetTypedAccount()
	if err != nil {
		return nil, errors.Base(err).Message("Shadowsocks|TCP: Failed to parse account.")
	}
	account := rawAccount.(*ShadowsocksAccount)

	iv := make([]byte, account.Cipher.IVSize())
	rand.Read(iv)
	_, err = writer.Write(iv)
	if err != nil {
		return nil, errors.Base(err).Message("Shadowsocks|TCP: Failed to write IV.")
	}

	stream, err := account.Cipher.NewEncodingStream(account.Key, iv)
	if err != nil {
		return nil, errors.Base(err).Message("Shadowsocks|TCP: Failed to create encoding stream.")
	}

	return v2io.NewAdaptiveWriter(crypto.NewCryptionWriter(stream, writer)), nil
}

func EncodeUDPPacket(request *protocol.RequestHeader, payload *alloc.Buffer) (*alloc.Buffer, error) {
	user := request.User
	rawAccount, err := user.GetTypedAccount()
	if err != nil {
		return nil, errors.Base(err).Message("Shadowsocks|UDP: Failed to parse account.")
	}
	account := rawAccount.(*ShadowsocksAccount)

	buffer := alloc.NewSmallBuffer()
	ivLen := account.Cipher.IVSize()
	buffer.Slice(0, ivLen)
	rand.Read(buffer.Value)
	iv := buffer.Value

	switch request.Address.Family() {
	case v2net.AddressFamilyIPv4:
		buffer.AppendBytes(AddrTypeIPv4)
		buffer.Append([]byte(request.Address.IP()))
	case v2net.AddressFamilyIPv6:
		buffer.AppendBytes(AddrTypeIPv6)
		buffer.Append([]byte(request.Address.IP()))
	case v2net.AddressFamilyDomain:
		buffer.AppendBytes(AddrTypeDomain, byte(len(request.Address.Domain())))
		buffer.Append([]byte(request.Address.Domain()))
	default:
		return nil, errors.New("Shadowsocks|UDP: Unsupported address type: ", request.Address.Family())
	}

	buffer.AppendUint16(uint16(request.Port))
	buffer.Append(payload.Value)

	if request.Option.Has(RequestOptionOneTimeAuth) {
		authenticator := NewAuthenticator(HeaderKeyGenerator(account.Key, iv))
		buffer.Value[ivLen] |= 0x10

		buffer.Value = authenticator.Authenticate(buffer.Value, buffer.Value[ivLen:])
	}

	stream, err := account.Cipher.NewEncodingStream(account.Key, iv)
	if err != nil {
		return nil, errors.Base(err).Message("Shadowsocks|TCP: Failed to create encoding stream.")
	}

	stream.XORKeyStream(buffer.Value[ivLen:], buffer.Value[ivLen:])
	return buffer, nil
}

func DecodeUDPPacket(user *protocol.User, payload *alloc.Buffer) (*protocol.RequestHeader, *alloc.Buffer, error) {
	rawAccount, err := user.GetTypedAccount()
	if err != nil {
		return nil, nil, errors.Base(err).Message("Shadowsocks|UDP: Failed to parse account.")
	}
	account := rawAccount.(*ShadowsocksAccount)

	ivLen := account.Cipher.IVSize()
	iv := payload.Value[:ivLen]
	payload.SliceFrom(ivLen)

	stream, err := account.Cipher.NewDecodingStream(account.Key, iv)
	if err != nil {
		return nil, nil, errors.Base(err).Message("Shadowsocks|UDP: Failed to initialize decoding stream.")
	}
	stream.XORKeyStream(payload.Value, payload.Value)

	authenticator := NewAuthenticator(HeaderKeyGenerator(account.Key, iv))
	request := &protocol.RequestHeader{
		Version: Version,
		User:    user,
		Command: protocol.RequestCommandUDP,
	}

	addrType := (payload.Value[0] & 0x0F)
	if (payload.Value[0] & 0x10) == 0x10 {
		request.Option |= RequestOptionOneTimeAuth
	}

	if request.Option.Has(RequestOptionOneTimeAuth) && account.OneTimeAuth == Account_Disabled {
		return nil, nil, errors.New("Shadowsocks|UDP: Rejecting packet with OTA enabled, while server disables OTA.")
	}

	if !request.Option.Has(RequestOptionOneTimeAuth) && account.OneTimeAuth == Account_Enabled {
		return nil, nil, errors.New("Shadowsocks|UDP: Rejecting packet with OTA disabled, while server enables OTA.")
	}

	if request.Option.Has(RequestOptionOneTimeAuth) {
		payloadLen := payload.Len() - AuthSize
		authBytes := payload.Value[payloadLen:]

		actualAuth := authenticator.Authenticate(nil, payload.Value[0:payloadLen])
		if !bytes.Equal(actualAuth, authBytes) {
			return nil, nil, errors.New("Shadowsocks|UDP: Invalid OTA.")
		}

		payload.Slice(0, payloadLen)
	}

	payload.SliceFrom(1)

	switch addrType {
	case AddrTypeIPv4:
		request.Address = v2net.IPAddress(payload.Value[:4])
		payload.SliceFrom(4)
	case AddrTypeIPv6:
		request.Address = v2net.IPAddress(payload.Value[:16])
		payload.SliceFrom(16)
	case AddrTypeDomain:
		domainLength := int(payload.Value[0])
		request.Address = v2net.DomainAddress(string(payload.Value[1 : 1+domainLength]))
		payload.SliceFrom(1 + domainLength)
	default:
		return nil, nil, errors.New("Shadowsocks|UDP: Unknown address type: ", addrType)
	}

	request.Port = v2net.PortFromBytes(payload.Value[:2])
	payload.SliceFrom(2)

	return request, payload, nil
}

type UDPReader struct {
	Reader io.Reader
	User   *protocol.User
}

func (v *UDPReader) Read() (*alloc.Buffer, error) {
	buffer := alloc.NewSmallBuffer()
	nBytes, err := v.Reader.Read(buffer.Value)
	if err != nil {
		buffer.Release()
		return nil, err
	}
	buffer.Slice(0, nBytes)
	_, payload, err := DecodeUDPPacket(v.User, buffer)
	if err != nil {
		buffer.Release()
		return nil, err
	}
	return payload, nil
}

func (v *UDPReader) Release() {
}

type UDPWriter struct {
	Writer  io.Writer
	Request *protocol.RequestHeader
}

func (v *UDPWriter) Write(buffer *alloc.Buffer) error {
	payload, err := EncodeUDPPacket(v.Request, buffer)
	if err != nil {
		return err
	}
	_, err = v.Writer.Write(payload.Value)
	payload.Release()
	return err
}

func (v *UDPWriter) Release() {

}
