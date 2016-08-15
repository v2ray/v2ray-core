package encoding

import (
	"crypto/md5"
	"crypto/rand"
	"hash/fnv"
	"io"

	"github.com/v2ray/v2ray-core/common/crypto"
	"github.com/v2ray/v2ray-core/common/log"
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/common/protocol"
	"github.com/v2ray/v2ray-core/proxy/vmess"
	"github.com/v2ray/v2ray-core/transport"
)

func hashTimestamp(t protocol.Timestamp) []byte {
	bytes := make([]byte, 0, 32)
	bytes = t.Bytes(bytes)
	bytes = t.Bytes(bytes)
	bytes = t.Bytes(bytes)
	bytes = t.Bytes(bytes)
	return bytes
}

type ClientSession struct {
	requestBodyKey  []byte
	requestBodyIV   []byte
	responseHeader  byte
	responseBodyKey []byte
	responseBodyIV  []byte
	responseReader  io.Reader
	idHash          protocol.IDHash
}

func NewClientSession(idHash protocol.IDHash) *ClientSession {
	randomBytes := make([]byte, 33) // 16 + 16 + 1
	rand.Read(randomBytes)

	session := &ClientSession{}
	session.requestBodyKey = randomBytes[:16]
	session.requestBodyIV = randomBytes[16:32]
	session.responseHeader = randomBytes[32]
	responseBodyKey := md5.Sum(session.requestBodyKey)
	responseBodyIV := md5.Sum(session.requestBodyIV)
	session.responseBodyKey = responseBodyKey[:]
	session.responseBodyIV = responseBodyIV[:]
	session.idHash = idHash

	return session
}

func (this *ClientSession) EncodeRequestHeader(header *protocol.RequestHeader, writer io.Writer) {
	timestamp := protocol.NewTimestampGenerator(protocol.NowTime(), 30)()
	idHash := this.idHash(header.User.Account.(*vmess.Account).AnyValidID().Bytes())
	idHash.Write(timestamp.Bytes(nil))
	writer.Write(idHash.Sum(nil))

	buffer := make([]byte, 0, 512)
	buffer = append(buffer, Version)
	buffer = append(buffer, this.requestBodyIV...)
	buffer = append(buffer, this.requestBodyKey...)
	buffer = append(buffer, this.responseHeader, byte(header.Option), byte(0), byte(0), byte(header.Command))
	buffer = header.Port.Bytes(buffer)

	switch header.Address.Family() {
	case v2net.AddressFamilyIPv4:
		buffer = append(buffer, AddrTypeIPv4)
		buffer = append(buffer, header.Address.IP()...)
	case v2net.AddressFamilyIPv6:
		buffer = append(buffer, AddrTypeIPv6)
		buffer = append(buffer, header.Address.IP()...)
	case v2net.AddressFamilyDomain:
		buffer = append(buffer, AddrTypeDomain, byte(len(header.Address.Domain())))
		buffer = append(buffer, header.Address.Domain()...)
	}

	fnv1a := fnv.New32a()
	fnv1a.Write(buffer)

	buffer = fnv1a.Sum(buffer)

	timestampHash := md5.New()
	timestampHash.Write(hashTimestamp(timestamp))
	iv := timestampHash.Sum(nil)
	account := header.User.Account.(*vmess.Account)
	aesStream := crypto.NewAesEncryptionStream(account.ID.CmdKey(), iv)
	aesStream.XORKeyStream(buffer, buffer)
	writer.Write(buffer)

	return
}

func (this *ClientSession) EncodeRequestBody(writer io.Writer) io.Writer {
	aesStream := crypto.NewAesEncryptionStream(this.requestBodyKey, this.requestBodyIV)
	return crypto.NewCryptionWriter(aesStream, writer)
}

func (this *ClientSession) DecodeResponseHeader(reader io.Reader) (*protocol.ResponseHeader, error) {
	aesStream := crypto.NewAesDecryptionStream(this.responseBodyKey, this.responseBodyIV)
	this.responseReader = crypto.NewCryptionReader(aesStream, reader)

	buffer := make([]byte, 256)

	_, err := io.ReadFull(this.responseReader, buffer[:4])
	if err != nil {
		log.Info("Raw: Failed to read response header: ", err)
		return nil, err
	}

	if buffer[0] != this.responseHeader {
		log.Info("Raw: Unexpected response header. Expecting ", this.responseHeader, " but actually ", buffer[0])
		return nil, transport.ErrCorruptedPacket
	}

	header := &protocol.ResponseHeader{
		Option: protocol.ResponseOption(buffer[1]),
	}

	if buffer[2] != 0 {
		cmdId := buffer[2]
		dataLen := int(buffer[3])
		_, err := io.ReadFull(this.responseReader, buffer[:dataLen])
		if err != nil {
			log.Info("Raw: Failed to read response command: ", err)
			return nil, err
		}
		data := buffer[:dataLen]
		command, err := UnmarshalCommand(cmdId, data)
		if err == nil {
			header.Command = command
		}
	}

	return header, nil
}

func (this *ClientSession) DecodeResponseBody(reader io.Reader) io.Reader {
	return this.responseReader
}
