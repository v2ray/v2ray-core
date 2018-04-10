package tls_test

import (
	gotls "crypto/tls"
	"crypto/x509"
	"testing"
	"time"

	"v2ray.com/core/common/protocol/tls/cert"
	. "v2ray.com/core/transport/internet/tls"
	. "v2ray.com/ext/assert"
)

func TestCertificateIssuing(t *testing.T) {
	assert := With(t)

	certificate := ParseCertificate(cert.MustGenerate(nil, cert.Authority(true), cert.KeyUsage(x509.KeyUsageCertSign)))
	certificate.Usage = Certificate_AUTHORITY_ISSUE

	c := &Config{
		Certificate: []*Certificate{
			certificate,
		},
	}

	tlsConfig := c.GetTLSConfig()
	v2rayCert, err := tlsConfig.GetCertificate(&gotls.ClientHelloInfo{
		ServerName: "www.v2ray.com",
	})
	assert(err, IsNil)

	x509Cert, err := x509.ParseCertificate(v2rayCert.Certificate[0])
	assert(err, IsNil)
	assert(x509Cert.NotAfter.After(time.Now()), IsTrue)
}
