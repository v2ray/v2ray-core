package raw

import (
	"crypto/md5"
	"crypto/rand"
	"hash/fnv"
	"io"

	"github.com/v2ray/v2ray-core/common/alloc"
	"github.com/v2ray/v2ray-core/common/crypto"
	"github.com/v2ray/v2ray-core/common/protocol"
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
	idHash          protocol.IDHash
}

func NewClientSession(idHash protocol.IDHash) *ClientSession {
	randomBytes := make([]byte, 33) // 16 + 16 + 1
	rand.Read(randomBytes)

	session := &ClientSession{}
	session.requestBodyKey = randomBytes[:16]
	session.requestBodyIV = randomBytes[16:32]
	session.responseHeader = randomBytes[32]
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

