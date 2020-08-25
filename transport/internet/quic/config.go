// +build !confonly

package quic

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"

	"golang.org/x/crypto/chacha20poly1305"
	"v2ray.com/core/common"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/transport/internet"
)

func getAuth(config *Config) (cipher.AEAD, error) {
	security := config.Security.GetSecurityType()
	if security == protocol.SecurityType_NONE {
		return nil, nil
	}

	salted := []byte(config.Key + "v2ray-quic-salt")
	key := sha256.Sum256(salted)

	if security == protocol.SecurityType_AES128_GCM {
		block, err := aes.NewCipher(key[:16])
		common.Must(err)
		return cipher.NewGCM(block)
	}

	if security == protocol.SecurityType_CHACHA20_POLY1305 {
		return chacha20poly1305.New(key[:])
	}

	return nil, newError("unsupported security type")
}

func getHeader(config *Config) (internet.PacketHeader, error) {
	if config.Header == nil {
		return nil, nil
	}

	msg, err := config.Header.GetInstance()
	if err != nil {
		return nil, err
	}

	return internet.CreatePacketHeader(msg)
}
