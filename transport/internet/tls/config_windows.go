// +build windows

package tls

import "crypto/x509"

func (c *Config) getCertPool() *x509.CertPool {
	return nil
}
