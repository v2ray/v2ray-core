package self_test

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"time"

	quic "github.com/lucas-clemente/quic-go"
	"github.com/lucas-clemente/quic-go/integrationtests/tools/testlog"
	"github.com/lucas-clemente/quic-go/integrationtests/tools/testserver"
	"github.com/lucas-clemente/quic-go/internal/protocol"
	"github.com/lucas-clemente/quic-go/internal/testdata"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Multiplexing", func() {
	for _, v := range append(protocol.SupportedVersions, protocol.VersionTLS) {
		version := v

		// gQUIC 44 uses 0 byte connection IDs for packets sent to the client
		// It's not possible to do demultiplexing.
		if v == protocol.Version44 {
			continue
		}

		Context(fmt.Sprintf("with QUIC version %s", version), func() {
			runServer := func(ln quic.Listener) {
				go func() {
					defer GinkgoRecover()
					for {
						sess, err := ln.Accept()
						if err != nil {
							return
						}
						go func() {
							defer GinkgoRecover()
							str, err := sess.OpenStream()
							Expect(err).ToNot(HaveOccurred())
							defer str.Close()
							_, err = str.Write(testserver.PRDataLong)
							Expect(err).ToNot(HaveOccurred())
						}()
					}
				}()
			}

			dial := func(conn net.PacketConn, addr net.Addr) {
				sess, err := quic.Dial(
					conn,
					addr,
					fmt.Sprintf("quic.clemente.io:%d", addr.(*net.UDPAddr).Port),
					nil,
					&quic.Config{Versions: []protocol.VersionNumber{version}},
				)
				Expect(err).ToNot(HaveOccurred())
				str, err := sess.AcceptStream()
				Expect(err).ToNot(HaveOccurred())
				data, err := ioutil.ReadAll(str)
				Expect(err).ToNot(HaveOccurred())
				Expect(data).To(Equal(testserver.PRDataLong))
			}

			Context("multiplexing clients on the same conn", func() {
				getListener := func() quic.Listener {
					ln, err := quic.ListenAddr(
						"localhost:0",
						testdata.GetTLSConfig(),
						&quic.Config{Versions: []protocol.VersionNumber{version}},
					)
					Expect(err).ToNot(HaveOccurred())
					return ln
				}

				It("multiplexes connections to the same server", func() {
					server := getListener()
					runServer(server)
					defer server.Close()

					addr, err := net.ResolveUDPAddr("udp", "localhost:0")
					Expect(err).ToNot(HaveOccurred())
					conn, err := net.ListenUDP("udp", addr)
					Expect(err).ToNot(HaveOccurred())
					defer conn.Close()

					done1 := make(chan struct{})
					done2 := make(chan struct{})
					go func() {
						defer GinkgoRecover()
						dial(conn, server.Addr())
						close(done1)
					}()
					go func() {
						defer GinkgoRecover()
						dial(conn, server.Addr())
						close(done2)
					}()
					timeout := 30 * time.Second
					if testlog.Debug() {
						timeout = time.Minute
					}
					Eventually(done1, timeout).Should(BeClosed())
					Eventually(done2, timeout).Should(BeClosed())
				})

				It("multiplexes connections to different servers", func() {
					server1 := getListener()
					runServer(server1)
					defer server1.Close()
					server2 := getListener()
					runServer(server2)
					defer server2.Close()

					addr, err := net.ResolveUDPAddr("udp", "localhost:0")
					Expect(err).ToNot(HaveOccurred())
					conn, err := net.ListenUDP("udp", addr)
					Expect(err).ToNot(HaveOccurred())
					defer conn.Close()

					done1 := make(chan struct{})
					done2 := make(chan struct{})
					go func() {
						defer GinkgoRecover()
						dial(conn, server1.Addr())
						close(done1)
					}()
					go func() {
						defer GinkgoRecover()
						dial(conn, server2.Addr())
						close(done2)
					}()
					timeout := 30 * time.Second
					if testlog.Debug() {
						timeout = time.Minute
					}
					Eventually(done1, timeout).Should(BeClosed())
					Eventually(done2, timeout).Should(BeClosed())
				})
			})

			Context("multiplexing server and client on the same conn", func() {
				It("connects to itself", func() {
					if version != protocol.VersionTLS {
						Skip("Connecting to itself only works with IETF QUIC.")
					}

					addr, err := net.ResolveUDPAddr("udp", "localhost:0")
					Expect(err).ToNot(HaveOccurred())
					conn, err := net.ListenUDP("udp", addr)
					Expect(err).ToNot(HaveOccurred())
					defer conn.Close()

					server, err := quic.Listen(
						conn,
						testdata.GetTLSConfig(),
						&quic.Config{Versions: []protocol.VersionNumber{version}},
					)
					Expect(err).ToNot(HaveOccurred())
					runServer(server)
					done := make(chan struct{})
					go func() {
						defer GinkgoRecover()
						dial(conn, server.Addr())
						close(done)
					}()
					timeout := 30 * time.Second
					if testlog.Debug() {
						timeout = time.Minute
					}
					Eventually(done, timeout).Should(BeClosed())
				})

				It("runs a server and client on the same conn", func() {
					if os.Getenv("CI") == "true" {
						Skip("This test is flaky on CIs, see see https://github.com/golang/go/issues/17677.")
					}
					addr1, err := net.ResolveUDPAddr("udp", "localhost:0")
					Expect(err).ToNot(HaveOccurred())
					conn1, err := net.ListenUDP("udp", addr1)
					Expect(err).ToNot(HaveOccurred())
					defer conn1.Close()

					addr2, err := net.ResolveUDPAddr("udp", "localhost:0")
					Expect(err).ToNot(HaveOccurred())
					conn2, err := net.ListenUDP("udp", addr2)
					Expect(err).ToNot(HaveOccurred())
					defer conn2.Close()

					server1, err := quic.Listen(
						conn1,
						testdata.GetTLSConfig(),
						&quic.Config{Versions: []protocol.VersionNumber{version}},
					)
					Expect(err).ToNot(HaveOccurred())
					runServer(server1)
					defer server1.Close()

					server2, err := quic.Listen(
						conn2,
						testdata.GetTLSConfig(),
						&quic.Config{Versions: []protocol.VersionNumber{version}},
					)
					Expect(err).ToNot(HaveOccurred())
					runServer(server2)
					defer server2.Close()

					done1 := make(chan struct{})
					done2 := make(chan struct{})
					go func() {
						defer GinkgoRecover()
						dial(conn2, server1.Addr())
						close(done1)
					}()
					go func() {
						defer GinkgoRecover()
						dial(conn1, server2.Addr())
						close(done2)
					}()
					timeout := 30 * time.Second
					if testlog.Debug() {
						timeout = time.Minute
					}
					Eventually(done1, timeout).Should(BeClosed())
					Eventually(done2, timeout).Should(BeClosed())
				})
			})
		})
	}
})
