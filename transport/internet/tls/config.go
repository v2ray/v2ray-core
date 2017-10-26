package tls

import (
	"crypto/tls"

	"v2ray.com/core/app/log"
)

var (
	globalSessionCache = tls.NewLRUClientSessionCache(128)
)

func (c *Config) BuildCertificates() []tls.Certificate {
	certs := make([]tls.Certificate, 0, len(c.Certificate))
	for _, entry := range c.Certificate {
		keyPair, err := tls.X509KeyPair(entry.Certificate, entry.Key)
		if err != nil {
			log.Trace(newError("ignoring invalid X509 key pair").Base(err).AtWarning())
			continue
		}
		certs = append(certs, keyPair)
	}
	return certs
}

func (c *Config) GetTLSConfig() *tls.Config {
	config := &tls.Config{
		ClientSessionCache: globalSessionCache,
		NextProtos:         []string{"http/1.1"},
	}
	if c == nil {
		return config
	}

	config.InsecureSkipVerify = c.AllowInsecure
	config.Certificates = c.BuildCertificates()
	config.BuildNameToCertificate()
	if len(c.ServerName) > 0 {
		config.ServerName = c.ServerName
	}

	return config
}

func (c *Config) OverrideServerNameIfEmpty(serverName string) {
	if len(c.ServerName) == 0 {
		c.ServerName = serverName
	}
}
