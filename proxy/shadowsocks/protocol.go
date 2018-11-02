package shadowsocks

import (
	"bytes"
	"crypto/rand"
	"io"

	"v2ray.com/core/common"
	"v2ray.com/core/common/bitmask"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/protocol"
)

const (
	Version                               = 1
	RequestOptionOneTimeAuth bitmask.Byte = 0x01
)

var addrParser = protocol.NewAddressParser(
	protocol.AddressFamilyByte(0x01, net.AddressFamilyIPv4),
	protocol.AddressFamilyByte(0x04, net.AddressFamilyIPv6),
	protocol.AddressFamilyByte(0x03, net.AddressFamilyDomain),
	protocol.WithAddressTypeParser(func(b byte) byte {
		return b & 0x0F
	}),
)

// ReadTCPSession reads a Shadowsocks TCP session from the given reader, returns its header and remaining parts.
func ReadTCPSession(user *protocol.MemoryUser, reader io.Reader) (*protocol.RequestHeader, buf.Reader, error) {
	account := user.Account.(*MemoryAccount)

	buffer := buf.New()
	defer buffer.Release()

	ivLen := account.Cipher.IVSize()
	var iv []byte
	if ivLen > 0 {
		if _, err := buffer.ReadFullFrom(reader, ivLen); err != nil {
			return nil, nil, newError("failed to read IV").Base(err)
		}

		iv = append([]byte(nil), buffer.BytesTo(ivLen)...)
	}

	r, err := account.Cipher.NewDecryptionReader(account.Key, iv, reader)
	if err != nil {
		return nil, nil, newError("failed to initialize decoding stream").Base(err).AtError()
	}
	br := &buf.BufferedReader{Reader: r}
	reader = nil

	authenticator := NewAuthenticator(HeaderKeyGenerator(account.Key, iv))
	request := &protocol.RequestHeader{
		Version: Version,
		User:    user,
		Command: protocol.RequestCommandTCP,
	}

	buffer.Clear()

	addr, port, err := addrParser.ReadAddressPort(buffer, br)
	if err != nil {
		return nil, nil, newError("failed to read address").Base(err)
	}

	request.Address = addr
	request.Port = port

	if !account.Cipher.IsAEAD() {
		if (buffer.Byte(0) & 0x10) == 0x10 {
			request.Option.Set(RequestOptionOneTimeAuth)
		}

		if request.Option.Has(RequestOptionOneTimeAuth) && account.OneTimeAuth == Account_Disabled {
			return nil, nil, newError("rejecting connection with OTA enabled, while server disables OTA")
		}

		if !request.Option.Has(RequestOptionOneTimeAuth) && account.OneTimeAuth == Account_Enabled {
			return nil, nil, newError("rejecting connection with OTA disabled, while server enables OTA")
		}
	}

	if request.Option.Has(RequestOptionOneTimeAuth) {
		actualAuth := make([]byte, AuthSize)
		authenticator.Authenticate(buffer.Bytes(), actualAuth)

		_, err := buffer.ReadFullFrom(br, AuthSize)
		if err != nil {
			return nil, nil, newError("Failed to read OTA").Base(err)
		}

		if !bytes.Equal(actualAuth, buffer.BytesFrom(-AuthSize)) {
			return nil, nil, newError("invalid OTA")
		}
	}

	if request.Address == nil {
		return nil, nil, newError("invalid remote address.")
	}

	var chunkReader buf.Reader
	if request.Option.Has(RequestOptionOneTimeAuth) {
		chunkReader = NewChunkReader(br, NewAuthenticator(ChunkKeyGenerator(iv)))
	} else {
		chunkReader = buf.NewReader(br)
	}

	return request, chunkReader, nil
}

// WriteTCPRequest writes Shadowsocks request into the given writer, and returns a writer for body.
func WriteTCPRequest(request *protocol.RequestHeader, writer io.Writer) (buf.Writer, error) {
	user := request.User
	account := user.Account.(*MemoryAccount)

	if account.Cipher.IsAEAD() {
		request.Option.Clear(RequestOptionOneTimeAuth)
	}

	var iv []byte
	if account.Cipher.IVSize() > 0 {
		iv = make([]byte, account.Cipher.IVSize())
		common.Must2(rand.Read(iv))
		if err := buf.WriteAllBytes(writer, iv); err != nil {
			return nil, newError("failed to write IV")
		}
	}

	w, err := account.Cipher.NewEncryptionWriter(account.Key, iv, writer)
	if err != nil {
		return nil, newError("failed to create encoding stream").Base(err).AtError()
	}

	header := buf.New()

	if err := addrParser.WriteAddressPort(header, request.Address, request.Port); err != nil {
		return nil, newError("failed to write address").Base(err)
	}

	if request.Option.Has(RequestOptionOneTimeAuth) {
		header.SetByte(0, header.Byte(0)|0x10)

		authenticator := NewAuthenticator(HeaderKeyGenerator(account.Key, iv))
		authBuffer := header.Extend(AuthSize)
		authenticator.Authenticate(header.Bytes(), authBuffer)
	}

	if err := w.WriteMultiBuffer(buf.NewMultiBufferValue(header)); err != nil {
		return nil, newError("failed to write header").Base(err)
	}

	var chunkWriter buf.Writer
	if request.Option.Has(RequestOptionOneTimeAuth) {
		chunkWriter = NewChunkWriter(w.(io.Writer), NewAuthenticator(ChunkKeyGenerator(iv)))
	} else {
		chunkWriter = w
	}

	return chunkWriter, nil
}

func ReadTCPResponse(user *protocol.MemoryUser, reader io.Reader) (buf.Reader, error) {
	account := user.Account.(*MemoryAccount)

	var iv []byte
	if account.Cipher.IVSize() > 0 {
		iv = make([]byte, account.Cipher.IVSize())
		if _, err := io.ReadFull(reader, iv); err != nil {
			return nil, newError("failed to read IV").Base(err)
		}
	}

	return account.Cipher.NewDecryptionReader(account.Key, iv, reader)
}

func WriteTCPResponse(request *protocol.RequestHeader, writer io.Writer) (buf.Writer, error) {
	user := request.User
	account := user.Account.(*MemoryAccount)

	var iv []byte
	if account.Cipher.IVSize() > 0 {
		iv = make([]byte, account.Cipher.IVSize())
		common.Must2(rand.Read(iv))
		if err := buf.WriteAllBytes(writer, iv); err != nil {
			return nil, newError("failed to write IV.").Base(err)
		}
	}

	return account.Cipher.NewEncryptionWriter(account.Key, iv, writer)
}

func EncodeUDPPacket(request *protocol.RequestHeader, payload []byte) (*buf.Buffer, error) {
	user := request.User
	account := user.Account.(*MemoryAccount)

	buffer := buf.New()
	ivLen := account.Cipher.IVSize()
	if ivLen > 0 {
		common.Must2(buffer.ReadFullFrom(rand.Reader, ivLen))
	}
	iv := buffer.Bytes()

	if err := addrParser.WriteAddressPort(buffer, request.Address, request.Port); err != nil {
		return nil, newError("failed to write address").Base(err)
	}

	buffer.Write(payload)

	if !account.Cipher.IsAEAD() && request.Option.Has(RequestOptionOneTimeAuth) {
		authenticator := NewAuthenticator(HeaderKeyGenerator(account.Key, iv))
		buffer.SetByte(ivLen, buffer.Byte(ivLen)|0x10)

		authPayload := buffer.BytesFrom(ivLen)
		authBuffer := buffer.Extend(AuthSize)
		authenticator.Authenticate(authPayload, authBuffer)
	}
	if err := account.Cipher.EncodePacket(account.Key, buffer); err != nil {
		return nil, newError("failed to encrypt UDP payload").Base(err)
	}

	return buffer, nil
}

func DecodeUDPPacket(user *protocol.MemoryUser, payload *buf.Buffer) (*protocol.RequestHeader, *buf.Buffer, error) {
	account := user.Account.(*MemoryAccount)

	var iv []byte
	if !account.Cipher.IsAEAD() && account.Cipher.IVSize() > 0 {
		// Keep track of IV as it gets removed from payload in DecodePacket.
		iv = make([]byte, account.Cipher.IVSize())
		copy(iv, payload.BytesTo(account.Cipher.IVSize()))
	}

	if err := account.Cipher.DecodePacket(account.Key, payload); err != nil {
		return nil, nil, newError("failed to decrypt UDP payload").Base(err)
	}

	request := &protocol.RequestHeader{
		Version: Version,
		User:    user,
		Command: protocol.RequestCommandUDP,
	}

	if !account.Cipher.IsAEAD() {
		if (payload.Byte(0) & 0x10) == 0x10 {
			request.Option |= RequestOptionOneTimeAuth
		}

		if request.Option.Has(RequestOptionOneTimeAuth) && account.OneTimeAuth == Account_Disabled {
			return nil, nil, newError("rejecting packet with OTA enabled, while server disables OTA").AtWarning()
		}

		if !request.Option.Has(RequestOptionOneTimeAuth) && account.OneTimeAuth == Account_Enabled {
			return nil, nil, newError("rejecting packet with OTA disabled, while server enables OTA").AtWarning()
		}

		if request.Option.Has(RequestOptionOneTimeAuth) {
			payloadLen := payload.Len() - AuthSize
			authBytes := payload.BytesFrom(payloadLen)

			authenticator := NewAuthenticator(HeaderKeyGenerator(account.Key, iv))
			actualAuth := make([]byte, AuthSize)
			authenticator.Authenticate(payload.BytesTo(payloadLen), actualAuth)
			if !bytes.Equal(actualAuth, authBytes) {
				return nil, nil, newError("invalid OTA")
			}

			payload.Resize(0, payloadLen)
		}
	}

	payload.SetByte(0, payload.Byte(0)&0x0F)

	addr, port, err := addrParser.ReadAddressPort(nil, payload)
	if err != nil {
		return nil, nil, newError("failed to parse address").Base(err)
	}

	request.Address = addr
	request.Port = port

	return request, payload, nil
}

type UDPReader struct {
	Reader io.Reader
	User   *protocol.MemoryUser
}

func (v *UDPReader) ReadMultiBuffer() (buf.MultiBuffer, error) {
	buffer := buf.New()
	_, err := buffer.ReadFrom(v.Reader)
	if err != nil {
		buffer.Release()
		return nil, err
	}
	_, payload, err := DecodeUDPPacket(v.User, buffer)
	if err != nil {
		buffer.Release()
		return nil, err
	}
	return buf.NewMultiBufferValue(payload), nil
}

type UDPWriter struct {
	Writer  io.Writer
	Request *protocol.RequestHeader
}

// Write implements io.Writer.
func (w *UDPWriter) Write(payload []byte) (int, error) {
	packet, err := EncodeUDPPacket(w.Request, payload)
	if err != nil {
		return 0, err
	}
	_, err = w.Writer.Write(packet.Bytes())
	packet.Release()
	return len(payload), err
}
