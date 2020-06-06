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
	authid := CreateAuthID(key[:], time.Now().Unix())

	nonce := make([]byte, 8)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}

	lengthbuf := bytes.NewBuffer(nil)

	var HeaderDataLen uint16
	HeaderDataLen = uint16(len(data))

	common.Must(binary.Write(lengthbuf, binary.BigEndian, HeaderDataLen))

	authidCheck := KDF16(key[:], "VMess AuthID Check Value", string(authid[:]), string(lengthbuf.Bytes()), string(nonce))

	lengthbufb := lengthbuf.Bytes()

	LengthMask := KDF16(key[:], "VMess AuthID Mask Value", string(authid[:]), string(nonce[:]))[:2]

	lengthbufb[0] = lengthbufb[0] ^ LengthMask[0]
	lengthbufb[1] = lengthbufb[1] ^ LengthMask[1]

	HeaderAEADKey := KDF16(key[:], "VMess Header AEAD Key", string(authid[:]), string(nonce))

	HeaderAEADNonce := KDF(key[:], "VMess Header AEAD Nonce", string(authid[:]), string(nonce))[:12]

	block, err := aes.NewCipher(HeaderAEADKey)
	if err != nil {
		panic(err.Error())
	}

	headerAEAD, err := cipher.NewGCM(block)

	if err != nil {
		panic(err.Error())
	}

	headerSealed := headerAEAD.Seal(nil, HeaderAEADNonce, data, authid[:])

	var outPutBuf = bytes.NewBuffer(nil)

	common.Must2(outPutBuf.Write(authid[:])) //16

	common.Must2(outPutBuf.Write(authidCheck)) //16

	common.Must2(outPutBuf.Write(lengthbufb)) //2

	common.Must2(outPutBuf.Write(nonce)) //8

	common.Must2(outPutBuf.Write(headerSealed))

	return outPutBuf.Bytes()
}

func OpenVMessAEADHeader(key [16]byte, authid [16]byte, data io.Reader) ([]byte, bool, error, int) {
	var authidCheck [16]byte
	var lengthbufb [2]byte
	var nonce [8]byte

	n, err := io.ReadFull(data, authidCheck[:])
	if err != nil {
		return nil, false, err, n
	}

	n2, err := io.ReadFull(data, lengthbufb[:])
	if err != nil {
		return nil, false, err, n + n2
	}

	n4, err := io.ReadFull(data, nonce[:])
	if err != nil {
		return nil, false, err, n + n2 + n4
	}

	//Unmask Length

	LengthMask := KDF16(key[:], "VMess AuthID Mask Value", string(authid[:]), string(nonce[:]))[:2]

	lengthbufb[0] = lengthbufb[0] ^ LengthMask[0]
	lengthbufb[1] = lengthbufb[1] ^ LengthMask[1]

	authidCheckV := KDF16(key[:], "VMess AuthID Check Value", string(authid[:]), string(lengthbufb[:]), string(nonce[:]))

	if !hmac.Equal(authidCheckV, authidCheck[:]) {
		return nil, true, errCheckMismatch, n + n2 + n4
	}

	var length uint16

	common.Must(binary.Read(bytes.NewReader(lengthbufb[:]), binary.BigEndian, &length))

	HeaderAEADKey := KDF16(key[:], "VMess Header AEAD Key", string(authid[:]), string(nonce[:]))

	HeaderAEADNonce := KDF(key[:], "VMess Header AEAD Nonce", string(authid[:]), string(nonce[:]))[:12]

	//16 == AEAD Tag size
	header := make([]byte, length+16)

	n3, err := io.ReadFull(data, header)
	if err != nil {
		return nil, false, err, n + n2 + n3 + n4
	}

	block, err := aes.NewCipher(HeaderAEADKey)
	if err != nil {
		panic(err.Error())
	}

	headerAEAD, err := cipher.NewGCM(block)

	if err != nil {
		panic(err.Error())
	}

	out, erropenAEAD := headerAEAD.Open(nil, HeaderAEADNonce, header, authid[:])

	if erropenAEAD != nil {
		return nil, true, erropenAEAD, n + n2 + n3 + n4
	}

	return out, false, nil, n + n2 + n3 + n4
}

var errCheckMismatch = errors.New("check verify failed")
