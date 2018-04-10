package tls

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"time"

	"v2ray.com/core/common/net"
	"v2ray.com/core/common/protocol/tls/cert"
	"v2ray.com/core/transport/internet"
)

var (
	globalSessionCache = tls.NewLRUClientSessionCache(128)
)

func ParseCertificate(c *cert.Certificate) *Certificate {
	certPEM, keyPEM := c.ToPEM()
	return &Certificate{
		Certificate: certPEM,
		Key:         keyPEM,
	}
}

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

func (c *Config) BuildCertificates() []tls.Certificate {
	certs := make([]tls.Certificate, 0, len(c.Certificate))
	for _, entry := range c.Certificate {
		if entry.Usage != Certificate_ENCIPHERMENT {
			continue
		}
		keyPair, err := tls.X509KeyPair(entry.Certificate, entry.Key)
		if err != nil {
			newError("ignoring invalid X509 key pair").Base(err).AtWarning().WriteToLog()
			continue
		}
		certs = append(certs, keyPair)
	}
	return certs
}

func isCertificateExpired(c *tls.Certificate) bool {
	// If leaf is not there, the certificate is probably not used yet. We trust user to provide a valid certificate.
	return c.Leaf != nil && c.Leaf.NotAfter.After(time.Now().Add(-time.Minute))
}

func issueCertificate(rawCA *Certificate, domain string) (*tls.Certificate, error) {
	parent, err := cert.ParseCertificate(rawCA.Certificate, rawCA.Key)
	if err != nil {
		return nil, newError("failed to parse raw certificate").Base(err)
	}
	newCert, err := cert.Generate(parent, cert.CommonName(domain), cert.DNSNames(domain))
	if err != nil {
		return nil, newError("failed to generate new certificate for ", domain).Base(err)
	}
	newCertPEM, newKeyPEM := newCert.ToPEM()
	cert, err := tls.X509KeyPair(newCertPEM, newKeyPEM)
	return &cert, err
}

func (c *Config) GetTLSConfig(opts ...Option) *tls.Config {
	config := &tls.Config{
		ClientSessionCache: globalSessionCache,
		RootCAs:            c.GetCertPool(),
	}
	if c == nil {
		return config
	}

	for _, opt := range opts {
		opt(config)
	}

	config.InsecureSkipVerify = c.AllowInsecure
	config.Certificates = c.BuildCertificates()
	config.BuildNameToCertificate()
	config.GetCertificate = func(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
		domain := hello.ServerName
		certExpired := false
		if certificate, found := config.NameToCertificate[domain]; found {
			if !isCertificateExpired(certificate) {
				return certificate, nil
			}
			certExpired = true
		}

		if certExpired {
			newCerts := make([]tls.Certificate, 0, len(config.Certificates))

			for _, certificate := range config.Certificates {
				if !isCertificateExpired(&certificate) {
					newCerts = append(newCerts, certificate)
				}
			}

			config.Certificates = newCerts
		}

		var issuedCertificate *tls.Certificate

		// Create a new certificate from existing CA if possible
		for _, rawCert := range c.Certificate {
			if rawCert.Usage == Certificate_AUTHORITY_ISSUE {
				newCert, err := issueCertificate(rawCert, domain)
				if err != nil {
					newError("failed to issue new certificate for ", domain).Base(err).WriteToLog()
					continue
				}

				config.Certificates = append(config.Certificates, *newCert)
				issuedCertificate = &config.Certificates[len(config.Certificates)-1]
				break
			}
		}

		if issuedCertificate == nil {
			return nil, newError("failed to create a new certificate for ", domain)
		}

		config.BuildNameToCertificate()

		return issuedCertificate, nil
	}
	if len(c.ServerName) > 0 {
		config.ServerName = c.ServerName
	}
	if len(c.NextProtocol) > 0 {
		config.NextProtos = c.NextProtocol
	}
	if len(config.NextProtos) == 0 {
		config.NextProtos = []string{"http/1.1"}
	}

	return config
}

type Option func(*tls.Config)

func WithDestination(dest net.Destination) Option {
	return func(config *tls.Config) {
		if dest.Address.Family().IsDomain() && len(config.ServerName) == 0 {
			config.ServerName = dest.Address.Domain()
		}
	}
}

func WithNextProto(protocol ...string) Option {
	return func(config *tls.Config) {
		if len(config.NextProtos) == 0 {
			config.NextProtos = protocol
		}
	}
}

func ConfigFromContext(ctx context.Context) *Config {
	securitySettings := internet.SecuritySettingsFromContext(ctx)
	if securitySettings == nil {
		return nil
	}
	config, ok := securitySettings.(*Config)
	if !ok {
		return nil
	}
	return config
}
