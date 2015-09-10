// Package vmess contains protocol definition, io lib for VMess.
package vmess

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	_ "log"
	mrand "math/rand"
	"net"

	"github.com/v2ray/v2ray-core"
	v2io "github.com/v2ray/v2ray-core/io"
	v2net "github.com/v2ray/v2ray-core/net"
)

const (
	addrTypeIPv4   = byte(0x01)
	addrTypeIPv6   = byte(0x03)
	addrTypeDomain = byte(0x02)

	Version = byte(0x01)
)

var (
	ErrorInvalidUser = errors.New("Invalid User")
)

// VMessRequest implements the request message of VMess protocol. It only contains
// the header of a request message. The data part will be handled by conection
// handler directly, in favor of data streaming.
// 1 Version
// 16 UserHash
// 16 Request IV
// 16 Request Key
// 4 Response Header
// 1 Command
// 2 Port
// 1 Address Type
// 256 Target Address

type VMessRequest [312]byte

func (r *VMessRequest) Version() byte {
	return r[0]
}

func (r *VMessRequest) SetVersion(version byte) *VMessRequest {
	r[0] = version
	return r
}

func (r *VMessRequest) UserHash() []byte {
	return r[1:17]
}

func (r *VMessRequest) RequestIV() []byte {
	return r[17:33]
}

func (r *VMessRequest) RequestKey() []byte {
	return r[33:49]
}

func (r *VMessRequest) ResponseHeader() []byte {
	return r[49:53]
}

func (r *VMessRequest) Command() byte {
	return r[53]
}

func (r *VMessRequest) SetCommand(command byte) *VMessRequest {
	r[53] = command
	return r
}

func (r *VMessRequest) Port() uint16 {
	return binary.BigEndian.Uint16(r.portBytes())
}

func (r *VMessRequest) portBytes() []byte {
	return r[54:56]
}

func (r *VMessRequest) SetPort(port uint16) *VMessRequest {
	binary.BigEndian.PutUint16(r.portBytes(), port)
	return r
}

func (r *VMessRequest) targetAddressType() byte {
	return r[56]
}

func (r *VMessRequest) Destination() v2net.VAddress {
	switch r.targetAddressType() {
	case addrTypeIPv4:
		fallthrough
	case addrTypeIPv6:
		return v2net.IPAddress(r.targetAddressBytes(), r.Port())
	case addrTypeDomain:
		return v2net.DomainAddress(r.TargetAddress(), r.Port())
	default:
		panic("Unpexected address type")
	}
}

func (r *VMessRequest) TargetAddress() string {
	switch r.targetAddressType() {
	case addrTypeIPv4:
		return net.IP(r[57:61]).String()
	case addrTypeIPv6:
		return net.IP(r[57:73]).String()
	case addrTypeDomain:
		domainLength := int(r[57])
		return string(r[58 : 58+domainLength])
	default:
		panic("Unexpected address type")
	}
}

func (r *VMessRequest) targetAddressBytes() []byte {
	switch r.targetAddressType() {
	case addrTypeIPv4:
		return r[57:61]
	case addrTypeIPv6:
		return r[57:73]
	case addrTypeDomain:
		domainLength := int(r[57])
		return r[57 : 58+domainLength]
	default:
		panic("Unexpected address type")
	}
}

func (r *VMessRequest) SetIPv4(ipv4 []byte) *VMessRequest {
	r[56] = addrTypeIPv4
	copy(r[57:], ipv4)
	return r
}

func (r *VMessRequest) SetIPv6(ipv6 []byte) *VMessRequest {
	r[56] = addrTypeIPv6
	copy(r[57:], ipv6)
	return r
}

func (r *VMessRequest) SetDomain(domain string) *VMessRequest {
	r[56] = addrTypeDomain
	r[57] = byte(len(domain))
	copy(r[58:], []byte(domain))
	return r
}

type VMessRequestReader struct {
	vUserSet *core.VUserSet
}

func NewVMessRequestReader(vUserSet *core.VUserSet) *VMessRequestReader {
	reader := new(VMessRequestReader)
	reader.vUserSet = vUserSet
	return reader
}

func (r *VMessRequestReader) Read(reader io.Reader) (*VMessRequest, error) {
	request := new(VMessRequest)

	nBytes, err := reader.Read(request[0:17] /* version + user hash */)
	if err != nil {
		return nil, err
	}
	if nBytes != 17 {
		err = fmt.Errorf("Unexpected length of header %d", nBytes)
		return nil, err
	}
	// TODO: verify version number
	userId, valid := r.vUserSet.IsValidUserId(request.UserHash())
	if !valid {
		return nil, ErrorInvalidUser
	}

	decryptor, err := NewDecryptionReader(reader, userId.Hash([]byte("PWD")), make([]byte, blockSize))
	if err != nil {
		return nil, err
	}

	buffer := make([]byte, 300)
	nBytes, err = decryptor.Read(buffer[0:1])
	if err != nil {
		return nil, err
	}

	randomLength := buffer[0]
	if randomLength <= 0 || randomLength > 32 {
		return nil, fmt.Errorf("Unexpected random length %d", randomLength)
	}
	_, err = decryptor.Read(buffer[:randomLength])
	if err != nil {
		return nil, err
	}

	// TODO: check number of bytes returned
	_, err = decryptor.Read(request.RequestIV())
	if err != nil {
		return nil, err
	}
	_, err = decryptor.Read(request.RequestKey())
	if err != nil {
		return nil, err
	}
	_, err = decryptor.Read(request.ResponseHeader())
	if err != nil {
		return nil, err
	}
	_, err = decryptor.Read(buffer[0:1])
	if err != nil {
		return nil, err
	}
	request.SetCommand(buffer[0])

	_, err = decryptor.Read(buffer[0:2])
	if err != nil {
		return nil, err
	}
	request.SetPort(binary.BigEndian.Uint16(buffer[0:2]))

	_, err = decryptor.Read(buffer[0:1])
	if err != nil {
		return nil, err
	}
	switch buffer[0] {
	case addrTypeIPv4:
		_, err = decryptor.Read(buffer[1:5])
		if err != nil {
			return nil, err
		}
		request.SetIPv4(buffer[1:5])
	case addrTypeIPv6:
		_, err = decryptor.Read(buffer[1:17])
		if err != nil {
			return nil, err
		}
		request.SetIPv6(buffer[1:17])
	case addrTypeDomain:
		_, err = decryptor.Read(buffer[1:2])
		if err != nil {
			return nil, err
		}
		domainLength := buffer[1]
		_, err = decryptor.Read(buffer[2 : 2+domainLength])
		if err != nil {
			return nil, err
		}
		request.SetDomain(string(buffer[2 : 2+domainLength]))
	}
	_, err = decryptor.Read(buffer[0:1])
	if err != nil {
		return nil, err
	}
	randomLength = buffer[0]
	_, err = decryptor.Read(buffer[:randomLength])
	if err != nil {
		return nil, err
	}

	return request, nil
}

type VMessRequestWriter struct {
	vUserSet *core.VUserSet
}

func NewVMessRequestWriter(vUserSet *core.VUserSet) *VMessRequestWriter {
	writer := new(VMessRequestWriter)
	writer.vUserSet = vUserSet
	return writer
}

func (w *VMessRequestWriter) Write(writer io.Writer, request *VMessRequest) error {
	buffer := make([]byte, 0, 300)
	buffer = append(buffer, request.Version())
	buffer = append(buffer, request.UserHash()...)

	encryptionBegin := len(buffer)

	randomLength := mrand.Intn(32) + 1
	randomContent := make([]byte, randomLength)
	_, err := rand.Read(randomContent)
	if err != nil {
		return err
	}
	buffer = append(buffer, byte(randomLength))
	buffer = append(buffer, randomContent...)

	buffer = append(buffer, request.RequestIV()...)
	buffer = append(buffer, request.RequestKey()...)
	buffer = append(buffer, request.ResponseHeader()...)
	buffer = append(buffer, request.Command())
	buffer = append(buffer, request.portBytes()...)
	buffer = append(buffer, request.targetAddressType())
	buffer = append(buffer, request.targetAddressBytes()...)

	paddingLength := blockSize - 1 - (len(buffer)-encryptionBegin)%blockSize
	if paddingLength == 0 {
		paddingLength = blockSize
	}
	paddingBuffer := make([]byte, paddingLength)
	_, err = rand.Read(paddingBuffer)
	if err != nil {
		return err
	}
	buffer = append(buffer, byte(paddingLength))
	buffer = append(buffer, paddingBuffer...)
	encryptionEnd := len(buffer)

	userId, valid := w.vUserSet.IsValidUserId(request.UserHash())
	if !valid {
		return ErrorInvalidUser
	}
	aesCipher, err := aes.NewCipher(userId.Hash([]byte("PWD")))
	if err != nil {
		return err
	}
	aesMode := cipher.NewCBCEncrypter(aesCipher, make([]byte, blockSize))
	cWriter := v2io.NewCryptionWriter(aesMode, writer)

	_, err = writer.Write(buffer[0:encryptionBegin])
	if err != nil {
		return err
	}
	_, err = cWriter.Write(buffer[encryptionBegin:encryptionEnd])
	if err != nil {
		return err
	}

	return nil
}

type VMessResponse [4]byte

func NewVMessResponse(request *VMessRequest) *VMessResponse {
	response := new(VMessResponse)
	copy(response[:], request.ResponseHeader())
	return response
}
