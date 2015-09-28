package protocol

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
	"hash/fnv"
	"time"

	"github.com/v2ray/v2ray-core/common/errors"
	"github.com/v2ray/v2ray-core/common/log"
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/proxy/vmess/protocol/user"
)

type VMessUDP struct {
	user    user.ID
	version byte
	token   uint16
	address v2net.Address
	data    []byte
}

func ReadVMessUDP(buffer []byte, userset user.UserSet) (*VMessUDP, error) {
	userHash := buffer[:user.IDBytesLen]
	userId, timeSec, valid := userset.GetUser(userHash)
	if !valid {
		return nil, errors.NewAuthenticationError(userHash)
	}

	buffer = buffer[user.IDBytesLen:]
	aesCipher, err := aes.NewCipher(userId.CmdKey())
	if err != nil {
		return nil, err
	}
	aesStream := cipher.NewCFBDecrypter(aesCipher, user.Int64Hash(timeSec))
	aesStream.XORKeyStream(buffer, buffer)

	fnvHash := binary.BigEndian.Uint32(buffer[:4])
	fnv1a := fnv.New32a()
	fnv1a.Write(buffer[4:])
	fnvHashActual := fnv1a.Sum32()

	if fnvHash != fnvHashActual {
		log.Warning("Unexpected fhv hash %d, should be %d", fnvHashActual, fnvHash)
		return nil, errors.NewCorruptedPacketError()
	}

	buffer = buffer[4:]

	vmess := &VMessUDP{
		user:    *userId,
		version: buffer[0],
		token:   binary.BigEndian.Uint16(buffer[1:3]),
	}

	// buffer[3] is reserved

	port := binary.BigEndian.Uint16(buffer[4:6])
	addrType := buffer[6]
	var address v2net.Address
	switch addrType {
	case addrTypeIPv4:
		address = v2net.IPAddress(buffer[7:11], port)
		buffer = buffer[11:]
	case addrTypeIPv6:
		address = v2net.IPAddress(buffer[7:23], port)
		buffer = buffer[23:]
	case addrTypeDomain:
		domainLength := buffer[7]
		domain := string(buffer[8 : 8+domainLength])
		address = v2net.DomainAddress(domain, port)
		buffer = buffer[8+domainLength:]
	default:
		log.Warning("Unexpected address type %d", addrType)
		return nil, errors.NewCorruptedPacketError()
	}

	vmess.address = address
	vmess.data = buffer

	return vmess, nil
}

func (vmess *VMessUDP) ToBytes(idHash user.CounterHash, randomRangeInt64 user.RandomInt64InRange, buffer []byte) []byte {
	if buffer == nil {
		buffer = make([]byte, 0, 2*1024)
	}

	counter := randomRangeInt64(time.Now().UTC().Unix(), 30)
	hash := idHash.Hash(vmess.user.Bytes[:], counter)

	buffer = append(buffer, hash...)
	encryptBegin := 16

	// Placeholder for fnv1a hash
	buffer = append(buffer, byte(0), byte(0), byte(0), byte(0))
	fnvHash := 16
	fnvHashBegin := 20

	buffer = append(buffer, vmess.version)
	buffer = append(buffer, byte(vmess.token>>8), byte(vmess.token))
	buffer = append(buffer, byte(0x00))
	buffer = append(buffer, vmess.address.PortBytes()...)
	switch {
	case vmess.address.IsIPv4():
		buffer = append(buffer, addrTypeIPv4)
		buffer = append(buffer, vmess.address.IP()...)
	case vmess.address.IsIPv6():
		buffer = append(buffer, addrTypeIPv6)
		buffer = append(buffer, vmess.address.IP()...)
	case vmess.address.IsDomain():
		buffer = append(buffer, addrTypeDomain)
		buffer = append(buffer, byte(len(vmess.address.Domain())))
		buffer = append(buffer, []byte(vmess.address.Domain())...)
	}

	buffer = append(buffer, vmess.data...)

	fnv1a := fnv.New32a()
	fnv1a.Write(buffer[fnvHashBegin:])
	fnvHashValue := fnv1a.Sum32()

	buffer[fnvHash] = byte(fnvHashValue >> 24)
	buffer[fnvHash+1] = byte(fnvHashValue >> 16)
	buffer[fnvHash+2] = byte(fnvHashValue >> 8)
	buffer[fnvHash+3] = byte(fnvHashValue)

	aesCipher, err := aes.NewCipher(vmess.user.CmdKey())
	if err != nil {
		log.Error("VMess failed to create AES cipher: %v", err)
		return nil
	}
	aesStream := cipher.NewCFBEncrypter(aesCipher, user.Int64Hash(counter))
	aesStream.XORKeyStream(buffer[encryptBegin:], buffer[encryptBegin:])

	return buffer
}
