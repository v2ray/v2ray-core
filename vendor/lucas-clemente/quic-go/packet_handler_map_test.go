package quic

import (
	"bytes"
	"errors"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/lucas-clemente/quic-go/internal/protocol"
	"github.com/lucas-clemente/quic-go/internal/utils"
	"github.com/lucas-clemente/quic-go/internal/wire"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Packet Handler Map", func() {
	var (
		handler *packetHandlerMap
		conn    *mockPacketConn
	)

	getPacket := func(connID protocol.ConnectionID) []byte {
		buf := &bytes.Buffer{}
		err := (&wire.Header{
			DestConnectionID: connID,
			PacketNumberLen:  protocol.PacketNumberLen1,
		}).Write(buf, protocol.PerspectiveServer, versionGQUICFrames)
		Expect(err).ToNot(HaveOccurred())
		return buf.Bytes()
	}

	BeforeEach(func() {
		conn = newMockPacketConn()
		handler = newPacketHandlerMap(conn, 5, utils.DefaultLogger).(*packetHandlerMap)
	})

	It("closes", func() {
		testErr := errors.New("test error	")
		sess1 := NewMockPacketHandler(mockCtrl)
		sess1.EXPECT().destroy(testErr)
		sess2 := NewMockPacketHandler(mockCtrl)
		sess2.EXPECT().destroy(testErr)
		handler.Add(protocol.ConnectionID{1, 1, 1, 1}, sess1)
		handler.Add(protocol.ConnectionID{2, 2, 2, 2}, sess2)
		handler.close(testErr)
	})

	Context("handling packets", func() {
		It("handles packets for different packet handlers on the same packet conn", func() {
			connID1 := protocol.ConnectionID{1, 2, 3, 4, 5, 6, 7, 8}
			connID2 := protocol.ConnectionID{8, 7, 6, 5, 4, 3, 2, 1}
			packetHandler1 := NewMockPacketHandler(mockCtrl)
			packetHandler2 := NewMockPacketHandler(mockCtrl)
			handledPacket1 := make(chan struct{})
			handledPacket2 := make(chan struct{})
			packetHandler1.EXPECT().handlePacket(gomock.Any()).Do(func(p *receivedPacket) {
				Expect(p.header.DestConnectionID).To(Equal(connID1))
				close(handledPacket1)
			})
			packetHandler1.EXPECT().GetVersion()
			packetHandler1.EXPECT().GetPerspective().Return(protocol.PerspectiveClient)
			packetHandler2.EXPECT().handlePacket(gomock.Any()).Do(func(p *receivedPacket) {
				Expect(p.header.DestConnectionID).To(Equal(connID2))
				close(handledPacket2)
			})
			packetHandler2.EXPECT().GetVersion()
			packetHandler2.EXPECT().GetPerspective().Return(protocol.PerspectiveClient)
			handler.Add(connID1, packetHandler1)
			handler.Add(connID2, packetHandler2)

			conn.dataToRead <- getPacket(connID1)
			conn.dataToRead <- getPacket(connID2)
			Eventually(handledPacket1).Should(BeClosed())
			Eventually(handledPacket2).Should(BeClosed())

			// makes the listen go routine return
			packetHandler1.EXPECT().destroy(gomock.Any()).AnyTimes()
			packetHandler2.EXPECT().destroy(gomock.Any()).AnyTimes()
			close(conn.dataToRead)
		})

		It("drops unparseable packets", func() {
			err := handler.handlePacket(nil, []byte("invalid"))
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("error parsing invariant header:"))
		})

		It("deletes nil session entries after a wait time", func() {
			handler.deleteClosedSessionsAfter = 10 * time.Millisecond
			connID := protocol.ConnectionID{1, 2, 3, 4, 5, 6, 7, 8}
			handler.Add(connID, NewMockPacketHandler(mockCtrl))
			handler.Remove(connID)
			Eventually(func() error {
				return handler.handlePacket(nil, getPacket(connID))
			}).Should(MatchError("received a packet with an unexpected connection ID 0x0102030405060708"))
		})

		It("ignores packets arriving late for closed sessions", func() {
			handler.deleteClosedSessionsAfter = time.Hour
			connID := protocol.ConnectionID{1, 2, 3, 4, 5, 6, 7, 8}
			handler.Add(connID, NewMockPacketHandler(mockCtrl))
			handler.Remove(connID)
			err := handler.handlePacket(nil, getPacket(connID))
			Expect(err).ToNot(HaveOccurred())
		})

		It("drops packets for unknown receivers", func() {
			connID := protocol.ConnectionID{1, 2, 3, 4, 5, 6, 7, 8}
			err := handler.handlePacket(nil, getPacket(connID))
			Expect(err).To(MatchError("received a packet with an unexpected connection ID 0x0102030405060708"))
		})

		It("errors on packets that are smaller than the Payload Length in the packet header", func() {
			connID := protocol.ConnectionID{1, 2, 3, 4, 5, 6, 7, 8}
			packetHandler := NewMockPacketHandler(mockCtrl)
			packetHandler.EXPECT().GetVersion().Return(versionIETFFrames)
			packetHandler.EXPECT().GetPerspective().Return(protocol.PerspectiveClient)
			handler.Add(connID, packetHandler)
			hdr := &wire.Header{
				IsLongHeader:     true,
				Type:             protocol.PacketTypeHandshake,
				PayloadLen:       1000,
				DestConnectionID: connID,
				PacketNumberLen:  protocol.PacketNumberLen1,
				Version:          versionIETFFrames,
			}
			buf := &bytes.Buffer{}
			Expect(hdr.Write(buf, protocol.PerspectiveServer, versionIETFFrames)).To(Succeed())
			buf.Write(bytes.Repeat([]byte{0}, 500))

			err := handler.handlePacket(nil, buf.Bytes())
			Expect(err).To(MatchError("packet payload (500 bytes) is smaller than the expected payload length (1000 bytes)"))
		})

		It("cuts packets at the Payload Length", func() {
			connID := protocol.ConnectionID{1, 2, 3, 4, 5, 6, 7, 8}
			packetHandler := NewMockPacketHandler(mockCtrl)
			packetHandler.EXPECT().GetVersion().Return(versionIETFFrames)
			packetHandler.EXPECT().GetPerspective().Return(protocol.PerspectiveClient)
			handler.Add(connID, packetHandler)
			packetHandler.EXPECT().handlePacket(gomock.Any()).Do(func(p *receivedPacket) {
				Expect(p.data).To(HaveLen(456))
			})

			hdr := &wire.Header{
				IsLongHeader:     true,
				Type:             protocol.PacketTypeHandshake,
				PayloadLen:       456,
				DestConnectionID: connID,
				PacketNumberLen:  protocol.PacketNumberLen1,
				Version:          versionIETFFrames,
			}
			buf := &bytes.Buffer{}
			Expect(hdr.Write(buf, protocol.PerspectiveServer, versionIETFFrames)).To(Succeed())
			buf.Write(bytes.Repeat([]byte{0}, 500))
			err := handler.handlePacket(nil, buf.Bytes())
			Expect(err).ToNot(HaveOccurred())
		})

		It("closes the packet handlers when reading from the conn fails", func() {
			done := make(chan struct{})
			packetHandler := NewMockPacketHandler(mockCtrl)
			packetHandler.EXPECT().destroy(gomock.Any()).Do(func(e error) {
				Expect(e).To(HaveOccurred())
				close(done)
			})
			handler.Add(protocol.ConnectionID{1, 2, 3, 4, 5, 6, 7, 8}, packetHandler)
			conn.Close()
			Eventually(done).Should(BeClosed())
		})
	})

	Context("running a server", func() {
		It("adds a server", func() {
			connID := protocol.ConnectionID{0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88}
			p := getPacket(connID)
			server := NewMockUnknownPacketHandler(mockCtrl)
			server.EXPECT().handlePacket(gomock.Any()).Do(func(p *receivedPacket) {
				Expect(p.header.DestConnectionID).To(Equal(connID))
			})
			handler.SetServer(server)
			Expect(handler.handlePacket(nil, p)).To(Succeed())
		})

		It("closes all server sessions", func() {
			clientSess := NewMockPacketHandler(mockCtrl)
			clientSess.EXPECT().GetPerspective().Return(protocol.PerspectiveClient)
			serverSess := NewMockPacketHandler(mockCtrl)
			serverSess.EXPECT().GetPerspective().Return(protocol.PerspectiveServer)
			serverSess.EXPECT().Close()

			handler.Add(protocol.ConnectionID{1, 1, 1, 1}, clientSess)
			handler.Add(protocol.ConnectionID{2, 2, 2, 2}, serverSess)
			handler.CloseServer()
		})

		It("stops handling packets with unknown connection IDs after the server is closed", func() {
			connID := protocol.ConnectionID{0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88}
			p := getPacket(connID)
			server := NewMockUnknownPacketHandler(mockCtrl)
			handler.SetServer(server)
			handler.CloseServer()
			Expect(handler.handlePacket(nil, p)).To(MatchError("received a packet with an unexpected connection ID 0x1122334455667788"))
		})
	})
})
