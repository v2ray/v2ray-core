// +build !confonly

package shadowsocks

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"hash"
	"hash/crc32"
	"io"
	"io/ioutil"
	"v2ray.com/core/common/dice"

	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/protocol"
)

const (
	Version = 1
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

	hashkdf := hmac.New(func() hash.Hash { return sha256.New() }, []byte("SSBSKDF"))
	hashkdf.Write(account.Key)

	behaviorSeed := crc32.ChecksumIEEE(hashkdf.Sum(nil))

	behaviorRand := dice.NewDeterministicDice(int64(behaviorSeed))
	BaseDrainSize := behaviorRand.Roll(3266)
	RandDrainMax := behaviorRand.Roll(64) + 1
	RandDrainRolled := dice.Roll(RandDrainMax)
	DrainSize := BaseDrainSize + 16 + 38 + RandDrainRolled
	readSizeRemain := DrainSize

	buffer := buf.New()
	defer buffer.Release()

	ivLen := account.Cipher.IVSize()
	var iv []byte
	if ivLen > 0 {
		if _, err := buffer.ReadFullFrom(reader, ivLen); err != nil {
			readSizeRemain -= int(buffer.Len())
			DrainConnN(reader, readSizeRemain)
			return nil, nil, newError("failed to read IV").Base(err)
		}

		iv = append([]byte(nil), buffer.BytesTo(ivLen)...)
	}

	r, err := account.Cipher.NewDecryptionReader(account.Key, iv, reader)
	if err != nil {
		readSizeRemain -= int(buffer.Len())
		DrainConnN(reader, readSizeRemain)
		return nil, nil, newError("failed to initialize decoding stream").Base(err).AtError()
	}
	br := &buf.BufferedReader{Reader: r}

	request := &protocol.RequestHeader{
		Version: Version,
		User:    user,
		Command: protocol.RequestCommandTCP,
	}

	readSizeRemain -= int(buffer.Len())
	buffer.Clear()

	addr, port, err := addrParser.ReadAddressPort(buffer, br)
	if err != nil {
		readSizeRemain -= int(buffer.Len())
		DrainConnN(reader, readSizeRemain)
		return nil, nil, newError("failed to read address").Base(err)
	}

	request.Address = addr
	request.Port = port

	if request.Address == nil {
		readSizeRemain -= int(buffer.Len())
		DrainConnN(reader, readSizeRemain)
		return nil, nil, newError("invalid remote address.")
	}

	return request, br, nil
}

func DrainConnN(reader io.Reader, n int) error {
	_, err := io.CopyN(ioutil.Discard, reader, int64(n))
	return err
}

// WriteTCPRequest writes Shadowsocks request into the given writer, and returns a writer for body.
func WriteTCPRequest(request *protocol.RequestHeader, writer io.Writer) (buf.Writer, error) {
	user := request.User
	account := user.Account.(*MemoryAccount)

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

	if err := w.WriteMultiBuffer(buf.MultiBuffer{header}); err != nil {
		return nil, newError("failed to write header").Base(err)
	}

	return w, nil
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

	if err := addrParser.WriteAddressPort(buffer, request.Address, request.Port); err != nil {
		return nil, newError("failed to write address").Base(err)
	}

	buffer.Write(payload)

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
	return buf.MultiBuffer{payload}, nil
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
