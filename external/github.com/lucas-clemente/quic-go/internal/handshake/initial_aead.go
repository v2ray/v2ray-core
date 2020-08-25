package handshake

import (
	"crypto"
	"crypto/aes"

	"v2ray.com/core/external/github.com/lucas-clemente/quic-go/internal/protocol"
	"v2ray.com/core/external/github.com/marten-seemann/qtls"
)

var quicVersion1Salt = []byte{0xef, 0x4f, 0xb0, 0xab, 0xb4, 0x74, 0x70, 0xc4, 0x1b, 0xef, 0xcf, 0x80, 0x31, 0x33, 0x4f, 0xae, 0x48, 0x5e, 0x09, 0xa0}

// NewInitialAEAD creates a new AEAD for Initial encryption / decryption.
func NewInitialAEAD(connID protocol.ConnectionID, pers protocol.Perspective) (Sealer, Opener, error) {
	clientSecret, serverSecret := computeSecrets(connID)
	var mySecret, otherSecret []byte
	if pers == protocol.PerspectiveClient {
		mySecret = clientSecret
		otherSecret = serverSecret
	} else {
		mySecret = serverSecret
		otherSecret = clientSecret
	}
	myKey, myHPKey, myIV := computeInitialKeyAndIV(mySecret)
	otherKey, otherHPKey, otherIV := computeInitialKeyAndIV(otherSecret)

	encrypter := qtls.AEADAESGCM13(myKey, myIV)
	hpEncrypter, err := aes.NewCipher(myHPKey)
	if err != nil {
		return nil, nil, err
	}
	decrypter := qtls.AEADAESGCM13(otherKey, otherIV)
	hpDecrypter, err := aes.NewCipher(otherHPKey)
	if err != nil {
		return nil, nil, err
	}
	return newSealer(encrypter, hpEncrypter, false), newOpener(decrypter, hpDecrypter, false), nil
}

func computeSecrets(connID protocol.ConnectionID) (clientSecret, serverSecret []byte) {
	initialSecret := qtls.HkdfExtract(crypto.SHA256, connID, quicVersion1Salt)
	clientSecret = qtls.HkdfExpandLabel(crypto.SHA256, initialSecret, []byte{}, "client in", crypto.SHA256.Size())
	serverSecret = qtls.HkdfExpandLabel(crypto.SHA256, initialSecret, []byte{}, "server in", crypto.SHA256.Size())
	return
}

func computeInitialKeyAndIV(secret []byte) (key, hpKey, iv []byte) {
	key = qtls.HkdfExpandLabel(crypto.SHA256, secret, []byte{}, "quic key", 16)
	hpKey = qtls.HkdfExpandLabel(crypto.SHA256, secret, []byte{}, "quic hp", 16)
	iv = qtls.HkdfExpandLabel(crypto.SHA256, secret, []byte{}, "quic iv", 12)
	return
}
