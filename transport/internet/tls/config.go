package tls

import (
	"context"
	"crypto/tls"

	"v2ray.com/core/common/net"
	"v2ray.com/core/transport/internet"
)

var (
	globalSessionCache = tls.NewLRUClientSessionCache(128)
)

func (c *Config) BuildCertificates() []tls.Certificate {
	certs := make([]tls.Certificate, 0, len(c.Certificate))
	for _, entry := range c.Certificate {
		keyPair, err := tls.X509KeyPair(entry.Certificate, entry.Key)
		if err != nil {
			newError("ignoring invalid X509 key pair").Base(err).AtWarning().WriteToLog()
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

type Option func(*Config)

func WithDestination(dest net.Destination) Option {
	return func(config *Config) {
		if dest.Address.Family().IsDomain() && len(config.ServerName) == 0 {
			config.ServerName = dest.Address.Domain()
		}
	}
}

func ConfigFromContext(ctx context.Context, opts ...Option) *Config {
	securitySettings := internet.SecuritySettingsFromContext(ctx)
	if securitySettings == nil {
		return nil
	}
	if config, ok := securitySettings.(*Config); ok {
		for _, opt := range opts {
			opt(config)
		}
		return config
	}
	return nil
}
