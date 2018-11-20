package quic

import (
	"crypto/tls"
	"crypto/x509"
	"errors"

	"github.com/lucas-clemente/quic-go/internal/protocol"
	"github.com/lucas-clemente/quic-go/internal/testdata"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Mint Utils", func() {
	Context("generating a mint.Config", func() {
		It("sets non-blocking mode", func() {
			mintConf, err := tlsToMintConfig(nil, protocol.PerspectiveClient)
			Expect(err).ToNot(HaveOccurred())
			Expect(mintConf.NonBlocking).To(BeTrue())
		})

		It("sets the certificate chain", func() {
			tlsConf := testdata.GetTLSConfig()
			mintConf, err := tlsToMintConfig(tlsConf, protocol.PerspectiveClient)
			Expect(err).ToNot(HaveOccurred())
			Expect(mintConf.Certificates).ToNot(BeEmpty())
			Expect(mintConf.Certificates).To(HaveLen(len(tlsConf.Certificates)))
		})

		It("copies values from the tls.Config", func() {
			verifyErr := errors.New("test err")
			certPool := &x509.CertPool{}
			tlsConf := &tls.Config{
				RootCAs:            certPool,
				ServerName:         "www.example.com",
				InsecureSkipVerify: true,
				VerifyPeerCertificate: func(_ [][]byte, _ [][]*x509.Certificate) error {
					return verifyErr
				},
			}
			mintConf, err := tlsToMintConfig(tlsConf, protocol.PerspectiveClient)
			Expect(err).ToNot(HaveOccurred())
			Expect(mintConf.RootCAs).To(Equal(certPool))
			Expect(mintConf.ServerName).To(Equal("www.example.com"))
			Expect(mintConf.InsecureSkipVerify).To(BeTrue())
			Expect(mintConf.VerifyPeerCertificate(nil, nil)).To(MatchError(verifyErr))
		})

		It("requires client authentication", func() {
			mintConf, err := tlsToMintConfig(nil, protocol.PerspectiveClient)
			Expect(err).ToNot(HaveOccurred())
			Expect(mintConf.RequireClientAuth).To(BeFalse())
			conf := &tls.Config{ClientAuth: tls.RequireAnyClientCert}
			mintConf, err = tlsToMintConfig(conf, protocol.PerspectiveClient)
			Expect(err).ToNot(HaveOccurred())
			Expect(mintConf.RequireClientAuth).To(BeTrue())
		})

		It("rejects unsupported client auth types", func() {
			conf := &tls.Config{ClientAuth: tls.RequireAndVerifyClientCert}
			_, err := tlsToMintConfig(conf, protocol.PerspectiveClient)
			Expect(err).To(MatchError("mint currently only support ClientAuthType RequireAnyClientCert"))
		})
	})
})
