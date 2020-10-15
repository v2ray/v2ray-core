package xtls_test

import (
	"crypto/x509"
	"testing"
	"time"

	xtls "github.com/xtls/go"

	"v2ray.com/core/common"
	"v2ray.com/core/common/protocol/tls/cert"
	. "v2ray.com/core/transport/internet/xtls"
)

func TestCertificateIssuing(t *testing.T) {
	certificate := ParseCertificate(cert.MustGenerate(nil, cert.Authority(true), cert.KeyUsage(x509.KeyUsageCertSign)))
	certificate.Usage = Certificate_AUTHORITY_ISSUE

	c := &Config{
		Certificate: []*Certificate{
			certificate,
		},
	}

	xtlsConfig := c.GetXTLSConfig()
	v2rayCert, err := xtlsConfig.GetCertificate(&xtls.ClientHelloInfo{
		ServerName: "www.v2fly.org",
	})
	common.Must(err)

	x509Cert, err := x509.ParseCertificate(v2rayCert.Certificate[0])
	common.Must(err)
	if !x509Cert.NotAfter.After(time.Now()) {
		t.Error("NotAfter: ", x509Cert.NotAfter)
	}
}

func TestExpiredCertificate(t *testing.T) {
	caCert := cert.MustGenerate(nil, cert.Authority(true), cert.KeyUsage(x509.KeyUsageCertSign))
	expiredCert := cert.MustGenerate(caCert, cert.NotAfter(time.Now().Add(time.Minute*-2)), cert.CommonName("www.v2fly.org"), cert.DNSNames("www.v2fly.org"))

	certificate := ParseCertificate(caCert)
	certificate.Usage = Certificate_AUTHORITY_ISSUE

	certificate2 := ParseCertificate(expiredCert)

	c := &Config{
		Certificate: []*Certificate{
			certificate,
			certificate2,
		},
	}

	xtlsConfig := c.GetXTLSConfig()
	v2rayCert, err := xtlsConfig.GetCertificate(&xtls.ClientHelloInfo{
		ServerName: "www.v2fly.org",
	})
	common.Must(err)

	x509Cert, err := x509.ParseCertificate(v2rayCert.Certificate[0])
	common.Must(err)
	if !x509Cert.NotAfter.After(time.Now()) {
		t.Error("NotAfter: ", x509Cert.NotAfter)
	}
}

func TestInsecureCertificates(t *testing.T) {
	c := &Config{
		AllowInsecureCiphers: true,
	}

	xtlsConfig := c.GetXTLSConfig()
	if len(xtlsConfig.CipherSuites) > 0 {
		t.Fatal("Unexpected tls cipher suites list: ", xtlsConfig.CipherSuites)
	}
}

func BenchmarkCertificateIssuing(b *testing.B) {
	certificate := ParseCertificate(cert.MustGenerate(nil, cert.Authority(true), cert.KeyUsage(x509.KeyUsageCertSign)))
	certificate.Usage = Certificate_AUTHORITY_ISSUE

	c := &Config{
		Certificate: []*Certificate{
			certificate,
		},
	}

	xtlsConfig := c.GetXTLSConfig()
	lenCerts := len(xtlsConfig.Certificates)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = xtlsConfig.GetCertificate(&xtls.ClientHelloInfo{
			ServerName: "www.v2fly.org",
		})
		delete(xtlsConfig.NameToCertificate, "www.v2fly.org")
		xtlsConfig.Certificates = xtlsConfig.Certificates[:lenCerts]
	}
}
