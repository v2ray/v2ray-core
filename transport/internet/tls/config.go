package tls

import (
	"crypto/tls"

	"v2ray.com/core/app/log"
	"v2ray.com/core/common/errors"
)

var (
	globalSessionCache = tls.NewLRUClientSessionCache(128)
)

func (v *Config) BuildCertificates() []tls.Certificate {
	certs := make([]tls.Certificate, 0, len(v.Certificate))
	for _, entry := range v.Certificate {
		keyPair, err := tls.X509KeyPair(entry.Certificate, entry.Key)
		if err != nil {
			log.Trace(errors.New("TLS: ignoring invalid X509 key pair").Base(err).AtWarning())
			continue
		}
		certs = append(certs, keyPair)
	}
	return certs
}

func (v *Config) GetTLSConfig() *tls.Config {
	config := &tls.Config{
		ClientSessionCache: globalSessionCache,
		NextProtos:         []string{"http/1.1"},
	}
	if v == nil {
		return config
	}

	config.InsecureSkipVerify = v.AllowInsecure
	config.Certificates = v.BuildCertificates()
	config.BuildNameToCertificate()
	if len(v.ServerName) > 0 {
		config.ServerName = v.ServerName
	}

	return config
}
