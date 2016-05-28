package http

import (
	"crypto/tls"

	v2net "github.com/v2ray/v2ray-core/common/net"
)

type CertificateConfig struct {
	Domain      string
	Certificate tls.Certificate
}

type TlsConfig struct {
	Enabled bool
	Certs   []*CertificateConfig
}

func (this *TlsConfig) GetConfig() *tls.Config {
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

type Config struct {
	OwnHosts  []v2net.Address
	TlsConfig *TlsConfig
}

func (this *Config) IsOwnHost(host v2net.Address) bool {
	for _, ownHost := range this.OwnHosts {
		if ownHost.Equals(host) {
			return true
		}
	}
	return false
}

type ClientConfig struct {
}
