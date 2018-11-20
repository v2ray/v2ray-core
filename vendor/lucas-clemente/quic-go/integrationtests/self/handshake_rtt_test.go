package self_test

import (
	"crypto/tls"
	"net"
	"time"

	quic "github.com/lucas-clemente/quic-go"
	"github.com/lucas-clemente/quic-go/integrationtests/tools/proxy"
	"github.com/lucas-clemente/quic-go/internal/protocol"
	"github.com/lucas-clemente/quic-go/qerr"

	"github.com/lucas-clemente/quic-go/internal/testdata"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Handshake RTT tests", func() {
	var (
		proxy         *quicproxy.QuicProxy
		server        quic.Listener
		serverConfig  *quic.Config
		testStartedAt time.Time
		acceptStopped chan struct{}
	)

	rtt := 400 * time.Millisecond

	BeforeEach(func() {
		acceptStopped = make(chan struct{})
		serverConfig = &quic.Config{}
	})

	AfterEach(func() {
		Expect(proxy.Close()).To(Succeed())
		Expect(server.Close()).To(Succeed())
		<-acceptStopped
	})

	runServerAndProxy := func() {
		var err error
		// start the server
		server, err = quic.ListenAddr("localhost:0", testdata.GetTLSConfig(), serverConfig)
		Expect(err).ToNot(HaveOccurred())
		// start the proxy
		proxy, err = quicproxy.NewQuicProxy("localhost:0", protocol.VersionWhatever, &quicproxy.Opts{
			RemoteAddr:  server.Addr().String(),
			DelayPacket: func(_ quicproxy.Direction, _ uint64) time.Duration { return rtt / 2 },
		})
		Expect(err).ToNot(HaveOccurred())

		testStartedAt = time.Now()

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

	expectDurationInRTTs := func(num int) {
		testDuration := time.Since(testStartedAt)
		rtts := float32(testDuration) / float32(rtt)
		Expect(rtts).To(SatisfyAll(
			BeNumerically(">=", num),
			BeNumerically("<", num+1),
		))
	}

	It("fails when there's no matching version, after 1 RTT", func() {
		if len(protocol.SupportedVersions) == 1 {
			Skip("Test requires at least 2 supported versions.")
		}
		serverConfig.Versions = protocol.SupportedVersions[:1]
		runServerAndProxy()
		clientConfig := &quic.Config{
			Versions: protocol.SupportedVersions[1:2],
		}
		_, err := quic.DialAddr(proxy.LocalAddr().String(), nil, clientConfig)
		Expect(err).To(HaveOccurred())
		Expect(err.(qerr.ErrorCode)).To(Equal(qerr.InvalidVersion))
		expectDurationInRTTs(1)
	})

	Context("gQUIC", func() {
		// 1 RTT for verifying the source address
		// 1 RTT to become secure
		// 1 RTT to become forward-secure
		It("is forward-secure after 3 RTTs", func() {
			runServerAndProxy()
			_, err := quic.DialAddr(proxy.LocalAddr().String(), &tls.Config{InsecureSkipVerify: true}, nil)
			Expect(err).ToNot(HaveOccurred())
			expectDurationInRTTs(3)
		})

		It("does version negotiation in 1 RTT, IETF QUIC => gQUIC", func() {
			clientConfig := &quic.Config{
				Versions: []protocol.VersionNumber{protocol.VersionTLS, protocol.SupportedVersions[0]},
			}
			runServerAndProxy()
			_, err := quic.DialAddr(
				proxy.LocalAddr().String(),
				&tls.Config{InsecureSkipVerify: true},
				clientConfig,
			)
			Expect(err).ToNot(HaveOccurred())
			expectDurationInRTTs(4)
		})

		It("is forward-secure after 2 RTTs when the server doesn't require a Cookie", func() {
			serverConfig.AcceptCookie = func(_ net.Addr, _ *quic.Cookie) bool {
				return true
			}
			runServerAndProxy()
			_, err := quic.DialAddr(proxy.LocalAddr().String(), &tls.Config{InsecureSkipVerify: true}, nil)
			Expect(err).ToNot(HaveOccurred())
			expectDurationInRTTs(2)
		})

		It("doesn't complete the handshake when the server never accepts the Cookie", func() {
			serverConfig.AcceptCookie = func(_ net.Addr, _ *quic.Cookie) bool {
				return false
			}
			runServerAndProxy()
			_, err := quic.DialAddr(proxy.LocalAddr().String(), &tls.Config{InsecureSkipVerify: true}, nil)
			Expect(err).To(HaveOccurred())
			Expect(err.(*qerr.QuicError).ErrorCode).To(Equal(qerr.CryptoTooManyRejects))
		})

		It("doesn't complete the handshake when the handshake timeout is too short", func() {
			serverConfig.HandshakeTimeout = 2 * rtt
			runServerAndProxy()
			_, err := quic.DialAddr(proxy.LocalAddr().String(), &tls.Config{InsecureSkipVerify: true}, nil)
			Expect(err).To(HaveOccurred())
			Expect(err.(*qerr.QuicError).ErrorCode).To(Equal(qerr.HandshakeTimeout))
			// 2 RTTs during the timeout
			// plus 1 RTT: the timer starts 0.5 RTTs after sending the first packet, and the CONNECTION_CLOSE needs another 0.5 RTTs to reach the client
			expectDurationInRTTs(3)
		})
	})

	Context("IETF QUIC", func() {
		var clientConfig *quic.Config
		var clientTLSConfig *tls.Config

		BeforeEach(func() {
			serverConfig.Versions = []protocol.VersionNumber{protocol.VersionTLS}
			clientConfig = &quic.Config{Versions: []protocol.VersionNumber{protocol.VersionTLS}}
			clientTLSConfig = &tls.Config{
				InsecureSkipVerify: true,
				ServerName:         "quic.clemente.io",
			}
		})

		// 1 RTT for verifying the source address
		// 1 RTT for the TLS handshake
		It("is forward-secure after 2 RTTs", func() {
			runServerAndProxy()
			_, err := quic.DialAddr(
				proxy.LocalAddr().String(),
				clientTLSConfig,
				clientConfig,
			)
			Expect(err).ToNot(HaveOccurred())
			expectDurationInRTTs(2)
		})

		It("does version negotiation in 1 RTT, gQUIC => IETF QUIC", func() {
			clientConfig.Versions = []protocol.VersionNumber{protocol.SupportedVersions[0], protocol.VersionTLS}
			runServerAndProxy()
			_, err := quic.DialAddr(
				proxy.LocalAddr().String(),
				clientTLSConfig,
				clientConfig,
			)
			Expect(err).ToNot(HaveOccurred())
			expectDurationInRTTs(3)
		})

		It("is forward-secure after 1 RTTs when the server doesn't require a Cookie", func() {
			serverConfig.AcceptCookie = func(_ net.Addr, _ *quic.Cookie) bool {
				return true
			}
			runServerAndProxy()
			_, err := quic.DialAddr(
				proxy.LocalAddr().String(),
				clientTLSConfig,
				clientConfig,
			)
			Expect(err).ToNot(HaveOccurred())
			expectDurationInRTTs(1)
		})

		It("doesn't complete the handshake when the server never accepts the Cookie", func() {
			serverConfig.AcceptCookie = func(_ net.Addr, _ *quic.Cookie) bool {
				return false
			}
			runServerAndProxy()
			_, err := quic.DialAddr(
				proxy.LocalAddr().String(),
				clientTLSConfig,
				clientConfig,
			)
			Expect(err).To(HaveOccurred())
			Expect(err.(qerr.ErrorCode)).To(Equal(qerr.CryptoTooManyRejects))
		})
	})
})
