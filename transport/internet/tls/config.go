package tls

import (
	"crypto/tls"

	"v2ray.com/core/common/log"
)

var (
	globalSessionCache = tls.NewLRUClientSessionCache(128)
)

func (this *Config) BuildCertificates() []tls.Certificate {
	certs := make([]tls.Certificate, 0, len(this.Certificate))
	for _, entry := range this.Certificate {
		keyPair, err := tls.X509KeyPair(entry.Certificate, entry.Key)
		if err != nil {
			log.Warning("TLS: ignoring invalid X509 key pair: ", err)
			continue
		}
		certs = append(certs, keyPair)
	}
	return certs
}

func (this *Config) GetTLSConfig() *tls.Config {
	config := &tls.Config{
		ClientSessionCache: globalSessionCache,
	}
	if this == nil {
		return config
	}

	config.InsecureSkipVerify = this.AllowInsecure
	config.Certificates = this.BuildCertificates()
	config.BuildNameToCertificate()

	return config
}
