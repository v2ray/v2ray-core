package testdata

import (
	"crypto/tls"
	"path"
	"runtime"
)

var certPath string

func init() {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("Failed to get current frame")
	}

	certPath = path.Join(path.Dir(path.Dir(path.Dir(filename))), "example")
}

// GetCertificatePaths returns the paths to 'fullchain.pem' and 'privkey.pem' for the
// quic.clemente.io cert.
func GetCertificatePaths() (string, string) {
	return path.Join(certPath, "fullchain.pem"), path.Join(certPath, "privkey.pem")
}

// GetTLSConfig returns a tls config for quic.clemente.io
func GetTLSConfig() *tls.Config {
	cert, err := tls.LoadX509KeyPair(GetCertificatePaths())
	if err != nil {
		panic(err)
	}
	return &tls.Config{
		Certificates: []tls.Certificate{cert},
	}
}

// GetCertificate returns a certificate for quic.clemente.io
func GetCertificate() tls.Certificate {
	cert, err := tls.LoadX509KeyPair(GetCertificatePaths())
	if err != nil {
		panic(err)
	}
	return cert
}
