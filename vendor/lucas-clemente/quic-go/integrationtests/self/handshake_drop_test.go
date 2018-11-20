package self_test

import (
	"fmt"
	mrand "math/rand"
	"net"
	"time"

	_ "github.com/lucas-clemente/quic-clients" // download clients
	"github.com/lucas-clemente/quic-go"
	"github.com/lucas-clemente/quic-go/integrationtests/tools/proxy"
	"github.com/lucas-clemente/quic-go/internal/protocol"
	"github.com/lucas-clemente/quic-go/internal/testdata"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var directions = []quicproxy.Direction{quicproxy.DirectionIncoming, quicproxy.DirectionOutgoing, quicproxy.DirectionBoth}

type applicationProtocol struct {
	name string
	run  func(protocol.VersionNumber)
}

var _ = Describe("Handshake drop tests", func() {
	var (
		proxy *quicproxy.QuicProxy
		ln    quic.Listener
	)

	startListenerAndProxy := func(dropCallback quicproxy.DropCallback, version protocol.VersionNumber) {
		var err error
		ln, err = quic.ListenAddr(
			"localhost:0",
			testdata.GetTLSConfig(),
			&quic.Config{
				Versions: []protocol.VersionNumber{version},
			},
		)
		Expect(err).ToNot(HaveOccurred())
		serverPort := ln.Addr().(*net.UDPAddr).Port
		proxy, err = quicproxy.NewQuicProxy("localhost:0", version, &quicproxy.Opts{
			RemoteAddr: fmt.Sprintf("localhost:%d", serverPort),
			DropPacket: dropCallback,
		},
		)
		Expect(err).ToNot(HaveOccurred())
	}

	stochasticDropper := func(freq int) bool {
		return mrand.Int63n(int64(freq)) == 0
	}

	clientSpeaksFirst := &applicationProtocol{
		name: "client speaks first",
		run: func(version protocol.VersionNumber) {
			serverSessionChan := make(chan quic.Session)
			go func() {
				defer GinkgoRecover()
				sess, err := ln.Accept()
				Expect(err).ToNot(HaveOccurred())
				defer sess.Close()
				str, err := sess.AcceptStream()
				Expect(err).ToNot(HaveOccurred())
				b := make([]byte, 6)
				_, err = gbytes.TimeoutReader(str, 10*time.Second).Read(b)
				Expect(err).ToNot(HaveOccurred())
				Expect(string(b)).To(Equal("foobar"))
				serverSessionChan <- sess
			}()
			sess, err := quic.DialAddr(
				fmt.Sprintf("quic.clemente.io:%d", proxy.LocalPort()),
				nil,
				&quic.Config{Versions: []protocol.VersionNumber{version}},
			)
			Expect(err).ToNot(HaveOccurred())
			str, err := sess.OpenStream()
			Expect(err).ToNot(HaveOccurred())
			_, err = str.Write([]byte("foobar"))
			Expect(err).ToNot(HaveOccurred())

			var serverSession quic.Session
			Eventually(serverSessionChan, 10*time.Second).Should(Receive(&serverSession))
			sess.Close()
			serverSession.Close()
		},
	}

	serverSpeaksFirst := &applicationProtocol{
		name: "server speaks first",
		run: func(version protocol.VersionNumber) {
			serverSessionChan := make(chan quic.Session)
			go func() {
				defer GinkgoRecover()
				sess, err := ln.Accept()
				Expect(err).ToNot(HaveOccurred())
				str, err := sess.OpenStream()
				Expect(err).ToNot(HaveOccurred())
				_, err = str.Write([]byte("foobar"))
				Expect(err).ToNot(HaveOccurred())
				serverSessionChan <- sess
			}()
			sess, err := quic.DialAddr(
				fmt.Sprintf("quic.clemente.io:%d", proxy.LocalPort()),
				nil,
				&quic.Config{Versions: []protocol.VersionNumber{version}},
			)
			Expect(err).ToNot(HaveOccurred())
			str, err := sess.AcceptStream()
			Expect(err).ToNot(HaveOccurred())
			b := make([]byte, 6)
			_, err = gbytes.TimeoutReader(str, 10*time.Second).Read(b)
			Expect(err).ToNot(HaveOccurred())
			Expect(string(b)).To(Equal("foobar"))

			var serverSession quic.Session
			Eventually(serverSessionChan, 10*time.Second).Should(Receive(&serverSession))
			sess.Close()
			serverSession.Close()
		},
	}

	nobodySpeaks := &applicationProtocol{
		name: "nobody speaks",
		run: func(version protocol.VersionNumber) {
			serverSessionChan := make(chan quic.Session)
			go func() {
				defer GinkgoRecover()
				sess, err := ln.Accept()
				Expect(err).ToNot(HaveOccurred())
				serverSessionChan <- sess
			}()
			sess, err := quic.DialAddr(
				fmt.Sprintf("quic.clemente.io:%d", proxy.LocalPort()),
				nil,
				&quic.Config{Versions: []protocol.VersionNumber{version}},
			)
			Expect(err).ToNot(HaveOccurred())
			var serverSession quic.Session
			Eventually(serverSessionChan, 10*time.Second).Should(Receive(&serverSession))
			// both server and client accepted a session. Close now.
			sess.Close()
			serverSession.Close()
		},
	}

	AfterEach(func() {
		Expect(proxy.Close()).To(Succeed())
	})

	for _, v := range append(protocol.SupportedVersions, protocol.VersionTLS) {
		version := v

		Context(fmt.Sprintf("with QUIC version %s", version), func() {
			for _, d := range directions {
				direction := d

				for _, a := range []*applicationProtocol{clientSpeaksFirst, serverSpeaksFirst, nobodySpeaks} {
					app := a

					Context(app.name, func() {
						It(fmt.Sprintf("establishes a connection when the first packet is lost in %s direction", d), func() {
							startListenerAndProxy(func(d quicproxy.Direction, p uint64) bool {
								return p == 1 && d.Is(direction)
							}, version)
							app.run(version)
						})

						It(fmt.Sprintf("establishes a connection when the second packet is lost in %s direction", d), func() {
							startListenerAndProxy(func(d quicproxy.Direction, p uint64) bool {
								return p == 2 && d.Is(direction)
							}, version)
							app.run(version)
						})

						It(fmt.Sprintf("establishes a connection when 1/5 of the packets are lost in %s direction", d), func() {
							startListenerAndProxy(func(d quicproxy.Direction, p uint64) bool {
								return d.Is(direction) && stochasticDropper(5)
							}, version)
							app.run(version)
						})
					})
				}
			}
		})
	}
})
