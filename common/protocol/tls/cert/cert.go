package cert

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"math/big"
	"time"
)

//go:generate go run $GOPATH/src/v2ray.com/core/common/errors/errorgen/main.go -pkg cert -path Protocol,TLS,Cert

type Certificate struct {
	// Cerificate in x509 format
	Certificate []byte
	// Private key in x509 format
	PrivateKey []byte
}

type Option func(*x509.Certificate)

func Authority(isCA bool) Option {
	return func(cert *x509.Certificate) {
		cert.IsCA = isCA
	}
}

func NotBefore(t time.Time) Option {
	return func(c *x509.Certificate) {
		c.NotBefore = t
	}
}

func NotAfter(t time.Time) Option {
	return func(c *x509.Certificate) {
		c.NotAfter = t
	}
}

func DNSNames(names ...string) Option {
	return func(c *x509.Certificate) {
		c.DNSNames = names
	}
}

func CommonName(name string) Option {
	return func(c *x509.Certificate) {
		c.Subject.CommonName = name
	}
}

func Generate(parent *x509.Certificate, opts ...Option) (*Certificate, error) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, newError("failed to generate RSA private key").Base(err)
	}

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, newError("failed to generate serial number").Base(err)
	}

	template := &x509.Certificate{
		SerialNumber:          serialNumber,
		NotBefore:             time.Now().Add(time.Hour * -1),
		NotAfter:              time.Now().Add(time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	for _, opt := range opts {
		opt(template)
	}

	if parent == nil {
		parent = template
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, template, parent, priv.Public(), priv)
	if err != nil {
		return nil, newError("failed to create certificate").Base(err)
	}

	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})

	return &Certificate{
		Certificate: certPEM,
		PrivateKey:  keyPEM,
	}, nil
}
