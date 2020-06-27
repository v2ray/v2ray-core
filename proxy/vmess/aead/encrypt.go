package aead

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/binary"
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

	aeadPayloadLengthSerializedByte := aeadPayloadLengthSerializeBuffer.Bytes()
	var payloadHeaderLengthAEADEncrypted []byte

	{
		payloadHeaderLengthAEADKey := KDF16(key[:], KDFSaltConst_VMessHeaderPayloadLengthAEADKey, string(generatedAuthID[:]), string(connectionNonce))

		payloadHeaderLengthAEADNonce := KDF(key[:], KDFSaltConst_VMessHeaderPayloadLengthAEADIV, string(generatedAuthID[:]), string(connectionNonce))[:12]

		payloadHeaderLengthAEADAESBlock, err := aes.NewCipher(payloadHeaderLengthAEADKey)
		if err != nil {
			panic(err.Error())
		}

		payloadHeaderAEAD, err := cipher.NewGCM(payloadHeaderLengthAEADAESBlock)

		if err != nil {
			panic(err.Error())
		}

		payloadHeaderLengthAEADEncrypted = payloadHeaderAEAD.Seal(nil, payloadHeaderLengthAEADNonce, aeadPayloadLengthSerializedByte, generatedAuthID[:])
	}

	var payloadHeaderAEADEncrypted []byte

	{
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

		payloadHeaderAEADEncrypted = payloadHeaderAEAD.Seal(nil, payloadHeaderAEADNonce, data, generatedAuthID[:])
	}

	var outputBuffer = bytes.NewBuffer(nil)

	common.Must2(outputBuffer.Write(generatedAuthID[:])) //16

	common.Must2(outputBuffer.Write(payloadHeaderLengthAEADEncrypted)) //2+16

	common.Must2(outputBuffer.Write(connectionNonce)) //8

	common.Must2(outputBuffer.Write(payloadHeaderAEADEncrypted))

	return outputBuffer.Bytes()
}

func OpenVMessAEADHeader(key [16]byte, authid [16]byte, data io.Reader) ([]byte, bool, error, int) {
	var payloadHeaderLengthAEADEncrypted [18]byte
	var nonce [8]byte

	var bytesRead int

	authidCheckValueReadBytesCounts, err := io.ReadFull(data, payloadHeaderLengthAEADEncrypted[:])
	bytesRead += authidCheckValueReadBytesCounts
	if err != nil {
		return nil, false, err, bytesRead
	}

	nonceReadBytesCounts, err := io.ReadFull(data, nonce[:])
	bytesRead += nonceReadBytesCounts
	if err != nil {
		return nil, false, err, bytesRead
	}

	//Decrypt Length

	var decryptedAEADHeaderLengthPayloadResult []byte

	{
		payloadHeaderLengthAEADKey := KDF16(key[:], KDFSaltConst_VMessHeaderPayloadLengthAEADKey, string(authid[:]), string(nonce[:]))

		payloadHeaderLengthAEADNonce := KDF(key[:], KDFSaltConst_VMessHeaderPayloadLengthAEADIV, string(authid[:]), string(nonce[:]))[:12]

		payloadHeaderAEADAESBlock, err := aes.NewCipher(payloadHeaderLengthAEADKey)
		if err != nil {
			panic(err.Error())
		}

		payloadHeaderLengthAEAD, err := cipher.NewGCM(payloadHeaderAEADAESBlock)

		if err != nil {
			panic(err.Error())
		}

		decryptedAEADHeaderLengthPayload, erropenAEAD := payloadHeaderLengthAEAD.Open(nil, payloadHeaderLengthAEADNonce, payloadHeaderLengthAEADEncrypted[:], authid[:])

		if erropenAEAD != nil {
			return nil, true, erropenAEAD, bytesRead
		}

		decryptedAEADHeaderLengthPayloadResult = decryptedAEADHeaderLengthPayload
	}

	var length uint16

	common.Must(binary.Read(bytes.NewReader(decryptedAEADHeaderLengthPayloadResult[:]), binary.BigEndian, &length))

	var decryptedAEADHeaderPayloadR []byte

	var payloadHeaderAEADEncryptedReadedBytesCounts int

	{
		payloadHeaderAEADKey := KDF16(key[:], KDFSaltConst_VMessHeaderPayloadAEADKey, string(authid[:]), string(nonce[:]))

		payloadHeaderAEADNonce := KDF(key[:], KDFSaltConst_VMessHeaderPayloadAEADIV, string(authid[:]), string(nonce[:]))[:12]

		//16 == AEAD Tag size
		payloadHeaderAEADEncrypted := make([]byte, length+16)

		payloadHeaderAEADEncryptedReadedBytesCounts, err = io.ReadFull(data, payloadHeaderAEADEncrypted)
		bytesRead += payloadHeaderAEADEncryptedReadedBytesCounts
		if err != nil {
			return nil, false, err, bytesRead
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
			return nil, true, erropenAEAD, bytesRead
		}

		decryptedAEADHeaderPayloadR = decryptedAEADHeaderPayload
	}

	return decryptedAEADHeaderPayloadR, false, nil, bytesRead
}
