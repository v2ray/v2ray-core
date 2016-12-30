package tls

import (
	"crypto/tls"

	"v2ray.com/core/common/log"
)

var (
	globalSessionCache = tls.NewLRUClientSessionCache(128)
)

func (v *Config) BuildCertificates() []tls.Certificate {
	certs := make([]tls.Certificate, 0, len(v.Certificate))
	for _, entry := range v.Certificate {
		keyPair, err := tls.X509KeyPair(entry.Certificate, entry.Key)
		if err != nil {
			log.Warning("TLS: ignoring invalid X509 key pair: ", err)
			continue
		}
		certs = append(certs, keyPair)
	}
	return certs
}

func (v *Config) GetTLSConfig() *tls.Config {
	config := &tls.Config{
		ClientSessionCache: globalSessionCache,
		//NextProtos:         []string{"http/2.0", "spdy/3", "http/1.1"},
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
