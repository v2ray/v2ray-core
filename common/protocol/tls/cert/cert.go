package cert

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/asn1"
	"encoding/pem"
	"math/big"
	"time"

	"v2ray.com/core/common"
)

//go:generate errorgen

type Certificate struct {
	// Cerificate in ASN.1 DER format
	Certificate []byte
	// Private key in ASN.1 DER format
	PrivateKey []byte
}

func ParseCertificate(certPEM []byte, keyPEM []byte) (*Certificate, error) {
	certBlock, _ := pem.Decode(certPEM)
	if certBlock == nil {
		return nil, newError("failed to decode certificate")
	}
	keyBlock, _ := pem.Decode(keyPEM)
	if keyBlock == nil {
		return nil, newError("failed to decode key")
	}
	return &Certificate{
		Certificate: certBlock.Bytes,
		PrivateKey:  keyBlock.Bytes,
	}, nil
}

func (c *Certificate) ToPEM() ([]byte, []byte) {
	return pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: c.Certificate}),
		pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: c.PrivateKey})
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

func KeyUsage(usage x509.KeyUsage) Option {
	return func(c *x509.Certificate) {
		c.KeyUsage = usage
	}
}

func Organization(org string) Option {
	return func(c *x509.Certificate) {
		c.Subject.Organization = []string{org}
	}
}

func MustGenerate(parent *Certificate, opts ...Option) *Certificate {
	cert, err := Generate(parent, opts...)
	common.Must(err)
	return cert
}

func publicKey(priv interface{}) interface{} {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &k.PublicKey
	case *ecdsa.PrivateKey:
		return &k.PublicKey
	case ed25519.PrivateKey:
		return k.Public().(ed25519.PublicKey)
	default:
		return nil
	}
}

func Generate(parent *Certificate, opts ...Option) (*Certificate, error) {
	var (
		pKey      interface{}
		parentKey interface{}
		err       error
	)
	// higher signing performance than RSA2048
	selfKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, newError("failed to generate self private key").Base(err)
	}
	parentKey = selfKey
	if parent != nil {
		if _, e := asn1.Unmarshal(parent.PrivateKey, &ecPrivateKey{}); e == nil {
			pKey, err = x509.ParseECPrivateKey(parent.PrivateKey)
		} else if _, e := asn1.Unmarshal(parent.PrivateKey, &pkcs8{}); e == nil {
			pKey, err = x509.ParsePKCS8PrivateKey(parent.PrivateKey)
		} else if _, e := asn1.Unmarshal(parent.PrivateKey, &pkcs1PrivateKey{}); e == nil {
			pKey, err = x509.ParsePKCS1PrivateKey(parent.PrivateKey)
		}
		if err != nil {
			return nil, newError("failed to parse parent private key").Base(err)
		}
		parentKey = pKey
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

	parentCert := template
	if parent != nil {
		pCert, err := x509.ParseCertificate(parent.Certificate)
		if err != nil {
			return nil, newError("failed to parse parent certificate").Base(err)
		}
		parentCert = pCert
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, template, parentCert, publicKey(selfKey), parentKey)
	if err != nil {
		return nil, newError("failed to create certificate").Base(err)
	}

	privateKey, err := x509.MarshalPKCS8PrivateKey(selfKey)
	if err != nil {
		return nil, newError("Unable to marshal private key").Base(err)
	}

	return &Certificate{
		Certificate: derBytes,
		PrivateKey:  privateKey,
	}, nil
}
