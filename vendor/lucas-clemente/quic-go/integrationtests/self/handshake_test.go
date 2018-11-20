package self_test

import (
	"crypto/tls"
	"fmt"
	"net"

	quic "github.com/lucas-clemente/quic-go"
	"github.com/lucas-clemente/quic-go/internal/protocol"
	"github.com/lucas-clemente/quic-go/internal/testdata"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type versioner interface {
	GetVersion() protocol.VersionNumber
}

var _ = Describe("Handshake tests", func() {
	var (
		server        quic.Listener
		serverConfig  *quic.Config
		acceptStopped chan struct{}
	)

	BeforeEach(func() {
		server = nil
		acceptStopped = make(chan struct{})
		serverConfig = &quic.Config{}
	})

	AfterEach(func() {
		if server != nil {
			server.Close()
			<-acceptStopped
		}
	})

	runServer := func() {
		var err error
		// start the server
		server, err = quic.ListenAddr("localhost:0", testdata.GetTLSConfig(), serverConfig)
		Expect(err).ToNot(HaveOccurred())

		go func() {
			defer GinkgoRecover()
			defer close(acceptStopped)
			for {
				_, err := server.Accept()
				if err != nil {
					return
				}
			}
		}()
	}

	Context("Version Negotiation", func() {
		var supportedVersions []protocol.VersionNumber

		BeforeEach(func() {
			supportedVersions = protocol.SupportedVersions
			protocol.SupportedVersions = append(protocol.SupportedVersions, []protocol.VersionNumber{7, 8, 9, 10}...)
		})

		AfterEach(func() {
			protocol.SupportedVersions = supportedVersions
		})

		It("when the server supports more versions than the client", func() {
			// the server doesn't support the highest supported version, which is the first one the client will try
			// but it supports a bunch of versions that the client doesn't speak
			serverConfig.Versions = []protocol.VersionNumber{protocol.SupportedVersions[1], 7, 8, 9}
			runServer()
			sess, err := quic.DialAddr(server.Addr().String(), &tls.Config{InsecureSkipVerify: true}, nil)
			Expect(err).ToNot(HaveOccurred())
			Expect(sess.(versioner).GetVersion()).To(Equal(protocol.SupportedVersions[1]))
		})

		It("when the client supports more versions than the server supports", func() {
			// the server doesn't support the highest supported version, which is the first one the client will try
			// but it supports a bunch of versions that the client doesn't speak
			serverConfig.Versions = supportedVersions
			runServer()
			conf := &quic.Config{
				Versions: []protocol.VersionNumber{7, 8, 9, protocol.SupportedVersions[1], 10},
			}
			sess, err := quic.DialAddr(server.Addr().String(), &tls.Config{InsecureSkipVerify: true}, conf)
			Expect(err).ToNot(HaveOccurred())
			Expect(sess.(versioner).GetVersion()).To(Equal(protocol.SupportedVersions[1]))
		})
	})

	Context("Certifiate validation", func() {
		for _, v := range []protocol.VersionNumber{protocol.Version39, protocol.VersionTLS} {
			version := v

			Context(fmt.Sprintf("using %s", version), func() {
				var clientConfig *quic.Config

				BeforeEach(func() {
					serverConfig.Versions = []protocol.VersionNumber{version}
					clientConfig = &quic.Config{
						Versions: []protocol.VersionNumber{version},
					}
				})

				It("accepts the certificate", func() {
					runServer()
					_, err := quic.DialAddr(fmt.Sprintf("quic.clemente.io:%d", server.Addr().(*net.UDPAddr).Port), nil, clientConfig)
					Expect(err).ToNot(HaveOccurred())
				})

				It("errors if the server name doesn't match", func() {
					runServer()
					_, err := quic.DialAddr(fmt.Sprintf("127.0.0.1:%d", server.Addr().(*net.UDPAddr).Port), nil, clientConfig)
					Expect(err).To(HaveOccurred())
				})

				It("uses the ServerName in the tls.Config", func() {
					runServer()
					conf := &tls.Config{ServerName: "quic.clemente.io"}
					_, err := quic.DialAddr(fmt.Sprintf("127.0.0.1:%d", server.Addr().(*net.UDPAddr).Port), conf, clientConfig)
					Expect(err).ToNot(HaveOccurred())
				})
			})
		}
	})
})
