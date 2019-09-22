// +build windows
// +build !confonly

package tls

import "crypto/x509"

func (c *Config) getCertPool() (*x509.CertPool, error) {
	return c.loadSelfCertPool()
}
