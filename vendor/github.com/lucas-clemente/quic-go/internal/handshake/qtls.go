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
	if c.MinVersion < qtls.VersionTLS13 {
		c.MinVersion = qtls.VersionTLS13
	}
	if c.MaxVersion < qtls.VersionTLS13 {
		c.MaxVersion = qtls.VersionTLS13
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
		MinVersion:                  c.MinVersion,
		MaxVersion:                  c.MaxVersion,
		CurvePreferences:            c.CurvePreferences,
		DynamicRecordSizingDisabled: c.DynamicRecordSizingDisabled,
		Renegotiation:               c.Renegotiation,
		KeyLogWriter:                c.KeyLogWriter,
	}
}
