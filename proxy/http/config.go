package http

import "crypto/tls"

// CertificateConfig is the config for TLS certificates used in HTTP proxy.
type CertificateConfig struct {
	Domain      string
	Certificate tls.Certificate
}

// TlsConfig is the config for TLS connections.
type TLSConfig struct {
	Enabled bool
	Certs   []*CertificateConfig
}

// GetConfig returns corresponding tls.Config.
func (this *TLSConfig) GetConfig() *tls.Config {
	if !this.Enabled {
		return nil
	}

	config := &tls.Config{
		InsecureSkipVerify: false,
	}

	config.Certificates = make([]tls.Certificate, len(this.Certs))
	for index, cert := range this.Certs {
		config.Certificates[index] = cert.Certificate
	}

	config.BuildNameToCertificate()

	return config
}

// Config for HTTP proxy server.
type Config struct {
	TLSConfig *TLSConfig
}

// ClientConfig for HTTP proxy client.
type ClientConfig struct {
}
