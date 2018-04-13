// +build !windows

package tls

import "crypto/x509"

func (c *Config) GetCertPool() *x509.CertPool {
	pool, err := x509.SystemCertPool()
	if err != nil {
		newError("failed to get system cert pool.").Base(err).WriteToLog()
		return nil
	}
	if pool != nil {
		for _, cert := range c.Certificate {
			if cert.Usage == Certificate_AUTHORITY_VERIFY {
				pool.AppendCertsFromPEM(cert.Certificate)
			}
		}
	}
	return pool
}
