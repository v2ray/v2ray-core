package handshake

import (
	"crypto"
	"crypto/aes"
	"crypto/cipher"

	"github.com/lucas-clemente/quic-go/internal/protocol"
	"github.com/marten-seemann/qtls"
)

var quicVersion1Salt = []byte{0xef, 0x4f, 0xb0, 0xab, 0xb4, 0x74, 0x70, 0xc4, 0x1b, 0xef, 0xcf, 0x80, 0x31, 0x33, 0x4f, 0xae, 0x48, 0x5e, 0x09, 0xa0}

func newInitialAEAD(connID protocol.ConnectionID, pers protocol.Perspective) (Sealer, Opener, error) {
	clientSecret, serverSecret := computeSecrets(connID)
	var mySecret, otherSecret []byte
	if pers == protocol.PerspectiveClient {
		mySecret = clientSecret
		otherSecret = serverSecret
	} else {
		mySecret = serverSecret
		otherSecret = clientSecret
	}
	myKey, myPNKey, myIV := computeInitialKeyAndIV(mySecret)
	otherKey, otherPNKey, otherIV := computeInitialKeyAndIV(otherSecret)

	encrypterCipher, err := aes.NewCipher(myKey)
	if err != nil {
		return nil, nil, err
	}
	encrypter, err := cipher.NewGCM(encrypterCipher)
	if err != nil {
		return nil, nil, err
	}
	pnEncrypter, err := aes.NewCipher(myPNKey)
	if err != nil {
		return nil, nil, err
	}
	decrypterCipher, err := aes.NewCipher(otherKey)
	if err != nil {
		return nil, nil, err
	}
	decrypter, err := cipher.NewGCM(decrypterCipher)
	if err != nil {
		return nil, nil, err
	}
	pnDecrypter, err := aes.NewCipher(otherPNKey)
	if err != nil {
		return nil, nil, err
	}
	return newSealer(encrypter, myIV, pnEncrypter, false), newOpener(decrypter, otherIV, pnDecrypter, false), nil
}

func computeSecrets(connID protocol.ConnectionID) (clientSecret, serverSecret []byte) {
	initialSecret := qtls.HkdfExtract(crypto.SHA256, connID, quicVersion1Salt)
	clientSecret = qtls.HkdfExpandLabel(crypto.SHA256, initialSecret, []byte{}, "client in", crypto.SHA256.Size())
	serverSecret = qtls.HkdfExpandLabel(crypto.SHA256, initialSecret, []byte{}, "server in", crypto.SHA256.Size())
	return
}

func computeInitialKeyAndIV(secret []byte) (key, pnKey, iv []byte) {
	key = qtls.HkdfExpandLabel(crypto.SHA256, secret, []byte{}, "quic key", 16)
	pnKey = qtls.HkdfExpandLabel(crypto.SHA256, secret, []byte{}, "quic hp", 16)
	iv = qtls.HkdfExpandLabel(crypto.SHA256, secret, []byte{}, "quic iv", 12)
	return
}
