package raw

import (
	"crypto/md5"
	"crypto/rand"
	"hash/fnv"
	"io"

	"github.com/v2ray/v2ray-core/common/alloc"
	"github.com/v2ray/v2ray-core/common/crypto"
	"github.com/v2ray/v2ray-core/common/log"
	"github.com/v2ray/v2ray-core/common/protocol"
	"github.com/v2ray/v2ray-core/transport"
)

func hashTimestamp(t protocol.Timestamp) []byte {
	once := t.Bytes()
	bytes := make([]byte, 0, 32)
	bytes = append(bytes, once...)
	bytes = append(bytes, once...)
	bytes = append(bytes, once...)
	bytes = append(bytes, once...)
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
	buffer := alloc.NewSmallBuffer().Clear()
	defer buffer.Release()

	timestamp := protocol.NewTimestampGenerator(protocol.NowTime(), 30)()
	idHash := this.idHash(header.User.AnyValidID().Bytes())
	idHash.Write(timestamp.Bytes())
	idHash.Sum(buffer.Value)

	encryptionBegin := buffer.Len()

	buffer.AppendBytes(Version)
	buffer.Append(this.requestBodyIV)
	buffer.Append(this.requestBodyKey)
	buffer.AppendBytes(this.responseHeader, byte(header.Option), byte(0), byte(0))
	buffer.AppendBytes(byte(header.Command))
	buffer.Append(header.Port.Bytes())

	switch {
	case header.Address.IsIPv4():
		buffer.AppendBytes(AddrTypeIPv4)
		buffer.Append(header.Address.IP())
	case header.Address.IsIPv6():
		buffer.AppendBytes(AddrTypeIPv6)
		buffer.Append(header.Address.IP())
	case header.Address.IsDomain():
		buffer.AppendBytes(AddrTypeDomain, byte(len(header.Address.Domain())))
		buffer.Append([]byte(header.Address.Domain()))
	}

	encryptionEnd := buffer.Len()

	fnv1a := fnv.New32a()
	fnv1a.Write(buffer.Value[encryptionBegin:encryptionEnd])

	fnvHash := fnv1a.Sum32()
	buffer.AppendBytes(byte(fnvHash>>24), byte(fnvHash>>16), byte(fnvHash>>8), byte(fnvHash))
	encryptionEnd += 4

	timestampHash := md5.New()
	timestampHash.Write(hashTimestamp(timestamp))
	iv := timestampHash.Sum(nil)
	aesStream := crypto.NewAesEncryptionStream(header.User.ID.CmdKey(), iv)
	aesStream.XORKeyStream(buffer.Value[encryptionBegin:encryptionEnd], buffer.Value[encryptionBegin:encryptionEnd])
	writer.Write(buffer.Value)

	return
}

func (this *ClientSession) EncodeRequestBody(writer io.Writer) io.Writer {
	aesStream := crypto.NewAesEncryptionStream(this.requestBodyKey, this.requestBodyIV)
	return crypto.NewCryptionWriter(aesStream, writer)
}

func (this *ClientSession) DecodeResponseHeader(reader io.Reader) (*protocol.ResponseHeader, error) {
	aesStream := crypto.NewAesDecryptionStream(this.responseBodyKey, this.responseBodyIV)
	this.responseReader = crypto.NewCryptionReader(aesStream, reader)

	buffer := alloc.NewSmallBuffer()
	defer buffer.Release()

	_, err := io.ReadFull(this.responseReader, buffer.Value[:4])
	if err != nil {
		log.Error("Raw: Failed to read response header: ", err)
		return nil, err
	}

	if buffer.Value[0] != this.responseHeader {
		log.Warning("Raw: Unexpected response header. Expecting %d, but actually %d", this.responseHeader, buffer.Value[0])
		return nil, transport.ErrorCorruptedPacket
	}

	header := new(protocol.ResponseHeader)

	if buffer.Value[2] != 0 {
		cmdId := buffer.Value[2]
		dataLen := int(buffer.Value[3])
		_, err := io.ReadFull(this.responseReader, buffer.Value[:dataLen])
		if err != nil {
			log.Error("Raw: Failed to read response command: ", err)
			return nil, err
		}
		data := buffer.Value[:dataLen]
		command, err := UnmarshalCommand(cmdId, data)
		header.Command = command
	}

	return header, nil
}

func (this *ClientSession) DecodeResponseBody(reader io.Reader) io.Reader {
	return this.responseReader
}
