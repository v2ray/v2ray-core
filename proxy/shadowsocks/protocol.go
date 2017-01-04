package shadowsocks

import (
	"bytes"
	"crypto/rand"
	"io"

	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/crypto"
	"v2ray.com/core/common/errors"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/common/serial"
)

const (
	Version                  = 1
	RequestOptionOneTimeAuth = protocol.RequestOption(101)

	AddrTypeIPv4   = 1
	AddrTypeIPv6   = 4
	AddrTypeDomain = 3
)

func ReadTCPSession(user *protocol.User, reader io.Reader) (*protocol.RequestHeader, buf.Reader, error) {
	rawAccount, err := user.GetTypedAccount()
	if err != nil {
		return nil, nil, errors.Base(err).Message("Shadowsocks|TCP: Failed to parse account.")
	}
	account := rawAccount.(*ShadowsocksAccount)

	buffer := buf.NewLocal(512)
	defer buffer.Release()

	ivLen := account.Cipher.IVSize()
	err = buffer.AppendSupplier(buf.ReadFullFrom(reader, ivLen))
	if err != nil {
		return nil, nil, errors.Base(err).Message("Shadowsocks|TCP: Failed to read IV.")
	}

	iv := append([]byte(nil), buffer.BytesTo(ivLen)...)

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
	err = buffer.AppendSupplier(buf.ReadFullFrom(reader, 1))
	if err != nil {
		return nil, nil, errors.Base(err).Message("Shadowsocks|TCP: Failed to read address type.")
	}

	addrType := (buffer.Byte(0) & 0x0F)
	if (buffer.Byte(0) & 0x10) == 0x10 {
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
		err := buffer.AppendSupplier(buf.ReadFullFrom(reader, 4))
		if err != nil {
			return nil, nil, errors.Base(err).Message("Shadowsocks|TCP: Failed to read IPv4 address.")
		}
		request.Address = v2net.IPAddress(buffer.BytesFrom(-4))
	case AddrTypeIPv6:
		err := buffer.AppendSupplier(buf.ReadFullFrom(reader, 16))
		if err != nil {
			return nil, nil, errors.Base(err).Message("Shadowsocks|TCP: Failed to read IPv6 address.")
		}
		request.Address = v2net.IPAddress(buffer.BytesFrom(-16))
	case AddrTypeDomain:
		err := buffer.AppendSupplier(buf.ReadFullFrom(reader, 1))
		if err != nil {
			return nil, nil, errors.Base(err).Message("Shadowsocks|TCP: Failed to read domain lenth.")
		}
		domainLength := int(buffer.BytesFrom(-1)[0])
		err = buffer.AppendSupplier(buf.ReadFullFrom(reader, domainLength))
		if err != nil {
			return nil, nil, errors.Base(err).Message("Shadowsocks|TCP: Failed to read domain.")
		}
		request.Address = v2net.DomainAddress(string(buffer.BytesFrom(-domainLength)))
	default:
		// Check address validity after OTA verification.
	}

	err = buffer.AppendSupplier(buf.ReadFullFrom(reader, 2))
	if err != nil {
		return nil, nil, errors.Base(err).Message("Shadowsocks|TCP: Failed to read port.")
	}
	request.Port = v2net.PortFromBytes(buffer.BytesFrom(-2))

	if request.Option.Has(RequestOptionOneTimeAuth) {
		actualAuth := make([]byte, AuthSize)
		authenticator.Authenticate(buffer.Bytes())(actualAuth)

		err := buffer.AppendSupplier(buf.ReadFullFrom(reader, AuthSize))
		if err != nil {
			return nil, nil, errors.Base(err).Message("Shadowsocks|TCP: Failed to read OTA.")
		}

		if !bytes.Equal(actualAuth, buffer.BytesFrom(-AuthSize)) {
			return nil, nil, errors.New("Shadowsocks|TCP: Invalid OTA")
		}
	}

	if request.Address == nil {
		return nil, nil, errors.New("Shadowsocks|TCP: Invalid remote address.")
	}

	var chunkReader buf.Reader
	if request.Option.Has(RequestOptionOneTimeAuth) {
		chunkReader = NewChunkReader(reader, NewAuthenticator(ChunkKeyGenerator(iv)))
	} else {
		chunkReader = buf.NewReader(reader)
	}

	return request, chunkReader, nil
}

func WriteTCPRequest(request *protocol.RequestHeader, writer io.Writer) (buf.Writer, error) {
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

	header := buf.NewLocal(512)

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

	header.AppendSupplier(serial.WriteUint16(uint16(request.Port)))

	if request.Option.Has(RequestOptionOneTimeAuth) {
		header.SetByte(0, header.Byte(0)|0x10)

		authenticator := NewAuthenticator(HeaderKeyGenerator(account.Key, iv))
		header.AppendSupplier(authenticator.Authenticate(header.Bytes()))
	}

	_, err = writer.Write(header.Bytes())
	if err != nil {
		return nil, errors.Base(err).Message("Shadowsocks|TCP: Failed to write header.")
	}

	var chunkWriter buf.Writer
	if request.Option.Has(RequestOptionOneTimeAuth) {
		chunkWriter = NewChunkWriter(writer, NewAuthenticator(ChunkKeyGenerator(iv)))
	} else {
		chunkWriter = buf.NewWriter(writer)
	}

	return chunkWriter, nil
}

func ReadTCPResponse(user *protocol.User, reader io.Reader) (buf.Reader, error) {
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
	return buf.NewReader(crypto.NewCryptionReader(stream, reader)), nil
}

func WriteTCPResponse(request *protocol.RequestHeader, writer io.Writer) (buf.Writer, error) {
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

	return buf.NewWriter(crypto.NewCryptionWriter(stream, writer)), nil
}

func EncodeUDPPacket(request *protocol.RequestHeader, payload *buf.Buffer) (*buf.Buffer, error) {
	user := request.User
	rawAccount, err := user.GetTypedAccount()
	if err != nil {
		return nil, errors.Base(err).Message("Shadowsocks|UDP: Failed to parse account.")
	}
	account := rawAccount.(*ShadowsocksAccount)

	buffer := buf.NewSmall()
	ivLen := account.Cipher.IVSize()
	buffer.AppendSupplier(buf.ReadFullFrom(rand.Reader, ivLen))
	iv := buffer.Bytes()

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

	buffer.AppendSupplier(serial.WriteUint16(uint16(request.Port)))
	buffer.Append(payload.Bytes())

	if request.Option.Has(RequestOptionOneTimeAuth) {
		authenticator := NewAuthenticator(HeaderKeyGenerator(account.Key, iv))
		buffer.SetByte(ivLen, buffer.Byte(ivLen)|0x10)

		buffer.AppendSupplier(authenticator.Authenticate(buffer.BytesFrom(ivLen)))
	}

	stream, err := account.Cipher.NewEncodingStream(account.Key, iv)
	if err != nil {
		return nil, errors.Base(err).Message("Shadowsocks|TCP: Failed to create encoding stream.")
	}

	stream.XORKeyStream(buffer.BytesFrom(ivLen), buffer.BytesFrom(ivLen))
	return buffer, nil
}

func DecodeUDPPacket(user *protocol.User, payload *buf.Buffer) (*protocol.RequestHeader, *buf.Buffer, error) {
	rawAccount, err := user.GetTypedAccount()
	if err != nil {
		return nil, nil, errors.Base(err).Message("Shadowsocks|UDP: Failed to parse account.")
	}
	account := rawAccount.(*ShadowsocksAccount)

	ivLen := account.Cipher.IVSize()
	iv := payload.BytesTo(ivLen)
	payload.SliceFrom(ivLen)

	stream, err := account.Cipher.NewDecodingStream(account.Key, iv)
	if err != nil {
		return nil, nil, errors.Base(err).Message("Shadowsocks|UDP: Failed to initialize decoding stream.")
	}
	stream.XORKeyStream(payload.Bytes(), payload.Bytes())

	authenticator := NewAuthenticator(HeaderKeyGenerator(account.Key, iv))
	request := &protocol.RequestHeader{
		Version: Version,
		User:    user,
		Command: protocol.RequestCommandUDP,
	}

	addrType := (payload.Byte(0) & 0x0F)
	if (payload.Byte(0) & 0x10) == 0x10 {
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
		authBytes := payload.BytesFrom(payloadLen)

		actualAuth := make([]byte, AuthSize)
		authenticator.Authenticate(payload.BytesTo(payloadLen))(actualAuth)
		if !bytes.Equal(actualAuth, authBytes) {
			return nil, nil, errors.New("Shadowsocks|UDP: Invalid OTA.")
		}

		payload.Slice(0, payloadLen)
	}

	payload.SliceFrom(1)

	switch addrType {
	case AddrTypeIPv4:
		request.Address = v2net.IPAddress(payload.BytesTo(4))
		payload.SliceFrom(4)
	case AddrTypeIPv6:
		request.Address = v2net.IPAddress(payload.BytesTo(16))
		payload.SliceFrom(16)
	case AddrTypeDomain:
		domainLength := int(payload.Byte(0))
		request.Address = v2net.DomainAddress(string(payload.BytesRange(1, 1+domainLength)))
		payload.SliceFrom(1 + domainLength)
	default:
		return nil, nil, errors.New("Shadowsocks|UDP: Unknown address type: ", addrType)
	}

	request.Port = v2net.PortFromBytes(payload.BytesTo(2))
	payload.SliceFrom(2)

	return request, payload, nil
}

type UDPReader struct {
	Reader io.Reader
	User   *protocol.User
}

func (v *UDPReader) Read() (*buf.Buffer, error) {
	buffer := buf.NewSmall()
	err := buffer.AppendSupplier(buf.ReadFrom(v.Reader))
	if err != nil {
		buffer.Release()
		return nil, err
	}
	_, payload, err := DecodeUDPPacket(v.User, buffer)
	if err != nil {
		buffer.Release()
		return nil, err
	}
	return payload, nil
}

type UDPWriter struct {
	Writer  io.Writer
	Request *protocol.RequestHeader
}

func (v *UDPWriter) Write(buffer *buf.Buffer) error {
	payload, err := EncodeUDPPacket(v.Request, buffer)
	if err != nil {
		return err
	}
	_, err = v.Writer.Write(payload.Bytes())
	payload.Release()
	return err
}
