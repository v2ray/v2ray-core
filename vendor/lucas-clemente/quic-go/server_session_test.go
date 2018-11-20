package quic

import (
	"github.com/lucas-clemente/quic-go/internal/protocol"
	"github.com/lucas-clemente/quic-go/internal/utils"
	"github.com/lucas-clemente/quic-go/internal/wire"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Server Session", func() {
	var (
		qsess *MockQuicSession
		sess  *serverSession
	)

	BeforeEach(func() {
		qsess = NewMockQuicSession(mockCtrl)
		sess = newServerSession(qsess, &Config{}, utils.DefaultLogger).(*serverSession)
	})

	It("handles packets", func() {
		p := &receivedPacket{
			header: &wire.Header{DestConnectionID: protocol.ConnectionID{1, 2, 3, 4, 5}},
		}
		qsess.EXPECT().handlePacket(p)
		sess.handlePacket(p)
	})

	It("ignores Public Resets", func() {
		p := &receivedPacket{
			header: &wire.Header{
				ResetFlag:        true,
				DestConnectionID: protocol.ConnectionID{0xde, 0xad, 0xbe, 0xef},
			},
		}
		err := sess.handlePacketImpl(p)
		Expect(err).To(MatchError("Received unexpected Public Reset for connection 0xdeadbeef"))
	})

	It("ignores delayed packets with mismatching versions, for gQUIC", func() {
		qsess.EXPECT().GetVersion().Return(protocol.VersionNumber(100))
		// don't EXPECT any calls to handlePacket()
		p := &receivedPacket{
			header: &wire.Header{
				VersionFlag:      true,
				Version:          protocol.VersionNumber(123),
				DestConnectionID: protocol.ConnectionID{0xde, 0xad, 0xbe, 0xef},
			},
		}
		err := sess.handlePacketImpl(p)
		Expect(err).ToNot(HaveOccurred())
	})

	It("ignores delayed packets with mismatching versions, for IETF QUIC", func() {
		qsess.EXPECT().GetVersion().Return(protocol.VersionNumber(100))
		// don't EXPECT any calls to handlePacket()
		p := &receivedPacket{
			header: &wire.Header{
				IsLongHeader:     true,
				Version:          protocol.VersionNumber(123),
				DestConnectionID: protocol.ConnectionID{0xde, 0xad, 0xbe, 0xef},
			},
		}
		err := sess.handlePacketImpl(p)
		Expect(err).ToNot(HaveOccurred())
	})

	It("ignores packets with the wrong Long Header type", func() {
		qsess.EXPECT().GetVersion().Return(protocol.VersionNumber(100))
		p := &receivedPacket{
			header: &wire.Header{
				IsLongHeader:     true,
				Type:             protocol.PacketTypeRetry,
				Version:          protocol.VersionNumber(100),
				DestConnectionID: protocol.ConnectionID{0xde, 0xad, 0xbe, 0xef},
			},
		}
		err := sess.handlePacketImpl(p)
		Expect(err).To(MatchError("Received unsupported packet type: Retry"))
	})

	It("passes on Handshake packets", func() {
		p := &receivedPacket{
			header: &wire.Header{
				IsLongHeader:     true,
				Type:             protocol.PacketTypeHandshake,
				Version:          protocol.VersionNumber(100),
				DestConnectionID: protocol.ConnectionID{0xde, 0xad, 0xbe, 0xef},
			},
		}
		qsess.EXPECT().GetVersion().Return(protocol.VersionNumber(100))
		qsess.EXPECT().handlePacket(p)
		Expect(sess.handlePacketImpl(p)).To(Succeed())
	})

	It("has the right perspective", func() {
		Expect(sess.GetPerspective()).To(Equal(protocol.PerspectiveServer))
	})
})
