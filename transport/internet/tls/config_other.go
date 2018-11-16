// +build !windows

package tls

import (
	"bytes"
	"crypto/x509"
	"sync"
)

type certPoolCache struct {
	sync.Mutex
	once       sync.Once
	pool       *x509.CertPool
	extraCerts [][]byte
}

func (c *certPoolCache) hasCert(cert []byte) bool {
	for _, xCert := range c.extraCerts {
		if bytes.Equal(xCert, cert) {
			return true
		}
	}
	return false
}

func (c *certPoolCache) get(extraCerts []*Certificate) *x509.CertPool {
	c.once.Do(func() {
		pool, err := x509.SystemCertPool()
		if err != nil {
			newError("failed to get system cert pool.").Base(err).WriteToLog()
			return
		}
		c.pool = pool
	})

	if c.pool == nil {
		return nil
	}

	if len(extraCerts) == 0 {
		return c.pool
	}

	c.Lock()
	defer c.Unlock()

	for _, cert := range extraCerts {
		if !c.hasCert(cert.Certificate) {
			c.pool.AppendCertsFromPEM(cert.Certificate)
			c.extraCerts = append(c.extraCerts, cert.Certificate)
		}
	}

	return c.pool
}

var combineCertPool certPoolCache

func (c *Config) getCertPool() *x509.CertPool {
	return combineCertPool.get(c.Certificate)
}
