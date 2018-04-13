// +build windows

package tls

import "crypto/x509"

func (c *Config) GetCertPool() *x509.CertPool {
	return nil
}
