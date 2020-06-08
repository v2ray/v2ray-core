package aead

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"io"
	"time"
	"v2ray.com/core/common"
)

func SealVMessAEADHeader(key [16]byte, data []byte) []byte {
	generatedAuthID := CreateAuthID(key[:], time.Now().Unix())

	connectionNonce := make([]byte, 8)
	if _, err := io.ReadFull(rand.Reader, connectionNonce); err != nil {
		panic(err.Error())
	}

	aeadPayloadLengthSerializeBuffer := bytes.NewBuffer(nil)

	var headerPayloadDataLen uint16
	headerPayloadDataLen = uint16(len(data))

	common.Must(binary.Write(aeadPayloadLengthSerializeBuffer, binary.BigEndian, headerPayloadDataLen))

	authidCheckValue := KDF16(key[:], KDFSaltConst_VmessAuthIDCheckValue, string(generatedAuthID[:]), string(aeadPayloadLengthSerializeBuffer.Bytes()), string(connectionNonce))

	aeadPayloadLengthSerializedByte := aeadPayloadLengthSerializeBuffer.Bytes()

	aeadPayloadLengthMask := KDF16(key[:], KDFSaltConst_VMessLengthMask, string(generatedAuthID[:]), string(connectionNonce[:]))[:2]

	aeadPayloadLengthSerializedByte[0] = aeadPayloadLengthSerializedByte[0] ^ aeadPayloadLengthMask[0]
	aeadPayloadLengthSerializedByte[1] = aeadPayloadLengthSerializedByte[1] ^ aeadPayloadLengthMask[1]

	payloadHeaderAEADKey := KDF16(key[:], KDFSaltConst_VMessHeaderPayloadAEADKey, string(generatedAuthID[:]), string(connectionNonce))

	payloadHeaderAEADNonce := KDF(key[:], KDFSaltConst_VMessHeaderPayloadAEADIV, string(generatedAuthID[:]), string(connectionNonce))[:12]

	payloadHeaderAEADAESBlock, err := aes.NewCipher(payloadHeaderAEADKey)
	if err != nil {
		panic(err.Error())
	}

	payloadHeaderAEAD, err := cipher.NewGCM(payloadHeaderAEADAESBlock)

	if err != nil {
		panic(err.Error())
	}

	payloadHeaderAEADEncrypted := payloadHeaderAEAD.Seal(nil, payloadHeaderAEADNonce, data, generatedAuthID[:])

	var outputBuffer = bytes.NewBuffer(nil)

	common.Must2(outputBuffer.Write(generatedAuthID[:])) //16

	common.Must2(outputBuffer.Write(authidCheckValue)) //16

	common.Must2(outputBuffer.Write(aeadPayloadLengthSerializedByte)) //2

	common.Must2(outputBuffer.Write(connectionNonce)) //8

	common.Must2(outputBuffer.Write(payloadHeaderAEADEncrypted))

	return outputBuffer.Bytes()
}

func OpenVMessAEADHeader(key [16]byte, authid [16]byte, data io.Reader) ([]byte, bool, error, int) {
	var authidCheckValue [16]byte
	var headerPayloadDataLen [2]byte
	var nonce [8]byte

	authidCheckValueReadBytesCounts, err := io.ReadFull(data, authidCheckValue[:])
	if err != nil {
		return nil, false, err, authidCheckValueReadBytesCounts
	}

	headerPayloadDataLenReadBytesCounts, err := io.ReadFull(data, headerPayloadDataLen[:])
	if err != nil {
		return nil, false, err, authidCheckValueReadBytesCounts + headerPayloadDataLenReadBytesCounts
	}

	nonceReadBytesCounts, err := io.ReadFull(data, nonce[:])
	if err != nil {
		return nil, false, err, authidCheckValueReadBytesCounts + headerPayloadDataLenReadBytesCounts + nonceReadBytesCounts
	}

	//Unmask Length

	LengthMask := KDF16(key[:], KDFSaltConst_VMessLengthMask, string(authid[:]), string(nonce[:]))[:2]

	headerPayloadDataLen[0] = headerPayloadDataLen[0] ^ LengthMask[0]
	headerPayloadDataLen[1] = headerPayloadDataLen[1] ^ LengthMask[1]

	authidCheckValueReceivedFromNetwork := KDF16(key[:], KDFSaltConst_VmessAuthIDCheckValue, string(authid[:]), string(headerPayloadDataLen[:]), string(nonce[:]))

	if !hmac.Equal(authidCheckValueReceivedFromNetwork, authidCheckValue[:]) {
		return nil, true, errCheckMismatch, authidCheckValueReadBytesCounts + headerPayloadDataLenReadBytesCounts + nonceReadBytesCounts
	}

	var length uint16

	common.Must(binary.Read(bytes.NewReader(headerPayloadDataLen[:]), binary.BigEndian, &length))

	payloadHeaderAEADKey := KDF16(key[:], KDFSaltConst_VMessHeaderPayloadAEADKey, string(authid[:]), string(nonce[:]))

	payloadHeaderAEADNonce := KDF(key[:], KDFSaltConst_VMessHeaderPayloadAEADIV, string(authid[:]), string(nonce[:]))[:12]

	//16 == AEAD Tag size
	payloadHeaderAEADEncrypted := make([]byte, length+16)

	payloadHeaderAEADEncryptedReadedBytesCounts, err := io.ReadFull(data, payloadHeaderAEADEncrypted)
	if err != nil {
		return nil, false, err, authidCheckValueReadBytesCounts + headerPayloadDataLenReadBytesCounts + payloadHeaderAEADEncryptedReadedBytesCounts + nonceReadBytesCounts
	}

	payloadHeaderAEADAESBlock, err := aes.NewCipher(payloadHeaderAEADKey)
	if err != nil {
		panic(err.Error())
	}

	payloadHeaderAEAD, err := cipher.NewGCM(payloadHeaderAEADAESBlock)

	if err != nil {
		panic(err.Error())
	}

	decryptedAEADHeaderPayload, erropenAEAD := payloadHeaderAEAD.Open(nil, payloadHeaderAEADNonce, payloadHeaderAEADEncrypted, authid[:])

	if erropenAEAD != nil {
		return nil, true, erropenAEAD, authidCheckValueReadBytesCounts + headerPayloadDataLenReadBytesCounts + payloadHeaderAEADEncryptedReadedBytesCounts + nonceReadBytesCounts
	}

	return decryptedAEADHeaderPayload, false, nil, authidCheckValueReadBytesCounts + headerPayloadDataLenReadBytesCounts + payloadHeaderAEADEncryptedReadedBytesCounts + nonceReadBytesCounts
}

var errCheckMismatch = errors.New("check verify failed")
