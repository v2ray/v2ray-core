package self

import (
	"fmt"
	"io/ioutil"
	"net"
	"time"

	_ "github.com/lucas-clemente/quic-clients" // download clients
	quic "github.com/lucas-clemente/quic-go"
	"github.com/lucas-clemente/quic-go/integrationtests/tools/proxy"
	"github.com/lucas-clemente/quic-go/integrationtests/tools/testserver"
	"github.com/lucas-clemente/quic-go/internal/protocol"
	"github.com/lucas-clemente/quic-go/internal/testdata"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("non-zero RTT", func() {
	for _, v := range append(protocol.SupportedVersions, protocol.VersionTLS) {
		version := v

		Context(fmt.Sprintf("with QUIC version %s", version), func() {
			roundTrips := [...]time.Duration{
				10 * time.Millisecond,
				50 * time.Millisecond,
				100 * time.Millisecond,
				200 * time.Millisecond,
			}

			for _, r := range roundTrips {
				rtt := r

				It(fmt.Sprintf("downloads a message with %s RTT", rtt), func() {
					ln, err := quic.ListenAddr(
						"localhost:0",
						testdata.GetTLSConfig(),
						&quic.Config{
							Versions: []protocol.VersionNumber{version},
						},
					)
					Expect(err).ToNot(HaveOccurred())
					done := make(chan struct{})
					go func() {
						defer GinkgoRecover()
						sess, err := ln.Accept()
						Expect(err).ToNot(HaveOccurred())
						str, err := sess.OpenStream()
						Expect(err).ToNot(HaveOccurred())
						_, err = str.Write(testserver.PRData)
						Expect(err).ToNot(HaveOccurred())
						str.Close()
						close(done)
					}()
					serverPort := ln.Addr().(*net.UDPAddr).Port
					proxy, err := quicproxy.NewQuicProxy("localhost:0", version, &quicproxy.Opts{
						RemoteAddr: fmt.Sprintf("localhost:%d", serverPort),
						DelayPacket: func(d quicproxy.Direction, p uint64) time.Duration {
							return rtt / 2
						},
					})
					Expect(err).ToNot(HaveOccurred())
					defer proxy.Close()

					sess, err := quic.DialAddr(
						fmt.Sprintf("quic.clemente.io:%d", proxy.LocalPort()),
						nil,
						&quic.Config{Versions: []protocol.VersionNumber{version}},
					)
					Expect(err).ToNot(HaveOccurred())
					str, err := sess.AcceptStream()
					Expect(err).ToNot(HaveOccurred())
					data, err := ioutil.ReadAll(str)
					Expect(err).ToNot(HaveOccurred())
					Expect(data).To(Equal(testserver.PRData))
					sess.Close()
					Eventually(done).Should(BeClosed())
				})
			}
		})
	}
})
