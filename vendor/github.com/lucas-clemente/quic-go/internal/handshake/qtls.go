package handshake

import (
	"crypto/tls"

	"github.com/marten-seemann/qtls"
)

func tlsConfigToQtlsConfig(c *tls.Config) *qtls.Config {
	if c == nil {
		c = &tls.Config{}
	}
	// QUIC requires TLS 1.3 or newer
	minVersion := c.MinVersion
	if minVersion < qtls.VersionTLS13 {
		minVersion = qtls.VersionTLS13
	}
	maxVersion := c.MaxVersion
	if maxVersion < qtls.VersionTLS13 {
		maxVersion = qtls.VersionTLS13
	}
	return &qtls.Config{
		Rand:              c.Rand,
		Time:              c.Time,
		Certificates:      c.Certificates,
		NameToCertificate: c.NameToCertificate,
		// TODO: make GetCertificate work
		// GetCertificate:              c.GetCertificate,
		GetClientCertificate: c.GetClientCertificate,
		// TODO: make GetConfigForClient work
		// GetConfigForClient:          c.GetConfigForClient,
		VerifyPeerCertificate:       c.VerifyPeerCertificate,
		RootCAs:                     c.RootCAs,
		NextProtos:                  c.NextProtos,
		ServerName:                  c.ServerName,
		ClientAuth:                  c.ClientAuth,
		ClientCAs:                   c.ClientCAs,
		InsecureSkipVerify:          c.InsecureSkipVerify,
		CipherSuites:                c.CipherSuites,
		PreferServerCipherSuites:    c.PreferServerCipherSuites,
		SessionTicketsDisabled:      c.SessionTicketsDisabled,
		SessionTicketKey:            c.SessionTicketKey,
		MinVersion:                  minVersion,
		MaxVersion:                  maxVersion,
		CurvePreferences:            c.CurvePreferences,
		DynamicRecordSizingDisabled: c.DynamicRecordSizingDisabled,
		Renegotiation:               c.Renegotiation,
		KeyLogWriter:                c.KeyLogWriter,
	}
}
