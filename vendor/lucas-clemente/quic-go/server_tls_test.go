package quic

import (
	"bytes"
	"net"

	"github.com/bifurcation/mint"
	"github.com/lucas-clemente/quic-go/internal/handshake"
	"github.com/lucas-clemente/quic-go/internal/protocol"
	"github.com/lucas-clemente/quic-go/internal/testdata"
	"github.com/lucas-clemente/quic-go/internal/utils"
	"github.com/lucas-clemente/quic-go/internal/wire"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Stateless TLS handling", func() {
	var (
		conn        *mockPacketConn
		server      *serverTLS
		sessionChan <-chan tlsSession
	)

	BeforeEach(func() {
		conn = newMockPacketConn()
		config := &Config{
			Versions: []protocol.VersionNumber{protocol.VersionTLS},
		}
		var err error
		server, sessionChan, err = newServerTLS(conn, config, nil, testdata.GetTLSConfig(), utils.DefaultLogger)
		Expect(err).ToNot(HaveOccurred())
	})

	parseHeader := func(data []byte) *wire.Header {
		b := bytes.NewReader(data)
		iHdr, err := wire.ParseInvariantHeader(b, 0)
		Expect(err).ToNot(HaveOccurred())
		hdr, err := iHdr.Parse(b, protocol.PerspectiveServer, protocol.VersionTLS)
		Expect(err).ToNot(HaveOccurred())
		return hdr
	}

	It("drops too small packets", func() {
		server.HandleInitial(&receivedPacket{
			header: &wire.Header{},
			data:   bytes.Repeat([]byte{0}, protocol.MinInitialPacketSize-1), // the packet is now 1 byte too small
		})
		Expect(conn.dataWritten.Len()).To(BeZero())
	})

	It("drops packets with a too short connection ID", func() {
		hdr := &wire.Header{
			SrcConnectionID:  protocol.ConnectionID{1, 2, 3, 4, 5, 6, 7, 8},
			DestConnectionID: protocol.ConnectionID{1, 2, 3, 4},
			PacketNumberLen:  protocol.PacketNumberLen1,
		}
		server.HandleInitial(&receivedPacket{
			header: hdr,
			data:   bytes.Repeat([]byte{0}, protocol.MinInitialPacketSize),
		})
		Expect(conn.dataWritten.Len()).To(BeZero())
	})

	It("replies with a Retry packet, if a Cookie is required", func() {
		server.config.AcceptCookie = func(_ net.Addr, _ *handshake.Cookie) bool { return false }
		hdr := &wire.Header{
			Type:             protocol.PacketTypeInitial,
			SrcConnectionID:  protocol.ConnectionID{5, 4, 3, 2, 1},
			DestConnectionID: protocol.ConnectionID{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			Version:          protocol.VersionTLS,
		}
		server.HandleInitial(&receivedPacket{
			remoteAddr: &net.UDPAddr{},
			header:     hdr,
			data:       bytes.Repeat([]byte{0}, protocol.MinInitialPacketSize),
		})
		Expect(conn.dataWritten.Len()).ToNot(BeZero())
		replyHdr := parseHeader(conn.dataWritten.Bytes())
		Expect(replyHdr.Type).To(Equal(protocol.PacketTypeRetry))
		Expect(replyHdr.SrcConnectionID).ToNot(Equal(hdr.DestConnectionID))
		Expect(replyHdr.SrcConnectionID.Len()).To(BeNumerically(">=", protocol.MinConnectionIDLenInitial))
		Expect(replyHdr.DestConnectionID).To(Equal(hdr.SrcConnectionID))
		Expect(replyHdr.OrigDestConnectionID).To(Equal(hdr.DestConnectionID))
		Expect(replyHdr.Token).ToNot(BeEmpty())
		Expect(sessionChan).ToNot(Receive())
	})

	It("creates a session, if no Cookie is required", func() {
		server.config.AcceptCookie = func(_ net.Addr, _ *handshake.Cookie) bool { return true }
		hdr := &wire.Header{
			Type:             protocol.PacketTypeInitial,
			SrcConnectionID:  protocol.ConnectionID{5, 4, 3, 2, 1},
			DestConnectionID: protocol.ConnectionID{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			Version:          protocol.VersionTLS,
		}
		p := &receivedPacket{
			header: hdr,
			data:   bytes.Repeat([]byte{0}, protocol.MinInitialPacketSize),
		}
		run := make(chan struct{})
		server.newSession = func(connection, sessionRunner, protocol.ConnectionID, protocol.ConnectionID, protocol.ConnectionID, protocol.PacketNumber, *Config, *mint.Config, *handshake.TransportParameters, utils.Logger, protocol.VersionNumber) (quicSession, error) {
			sess := NewMockQuicSession(mockCtrl)
			sess.EXPECT().handlePacket(p)
			sess.EXPECT().run().Do(func() { close(run) })
			return sess, nil
		}

		done := make(chan struct{})
		go func() {
			defer GinkgoRecover()
			server.HandleInitial(p)
			// the Handshake packet is written by the session
			Expect(conn.dataWritten.Len()).To(BeZero())
			close(done)
		}()
		var tlsSess tlsSession
		Eventually(sessionChan).Should(Receive(&tlsSess))
		// make sure we're using a server-generated connection ID
		Expect(tlsSess.connID).ToNot(Equal(hdr.SrcConnectionID))
		Expect(tlsSess.connID).ToNot(Equal(hdr.DestConnectionID))
		Eventually(run).Should(BeClosed())
		Eventually(done).Should(BeClosed())
	})
})
