package quic

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	"v2ray.com/core/external/github.com/lucas-clemente/quic-go/internal/protocol"
	"v2ray.com/core/external/github.com/lucas-clemente/quic-go/internal/utils"
	"v2ray.com/core/external/github.com/lucas-clemente/quic-go/internal/wire"
)

type packetHandlerEntry struct {
	handler    packetHandler
	resetToken *[16]byte
}

// The packetHandlerMap stores packetHandlers, identified by connection ID.
// It is used:
// * by the server to store sessions
// * when multiplexing outgoing connections to store clients
type packetHandlerMap struct {
	mutex sync.RWMutex

	conn      net.PacketConn
	connIDLen int

	handlers    map[string] /* string(ConnectionID)*/ packetHandlerEntry
	resetTokens map[[16]byte] /* stateless reset token */ packetHandler
	server      unknownPacketHandler
	closed      bool

	deleteRetiredSessionsAfter time.Duration

	logger utils.Logger
}

var _ packetHandlerManager = &packetHandlerMap{}

func newPacketHandlerMap(conn net.PacketConn, connIDLen int, logger utils.Logger) packetHandlerManager {
	m := &packetHandlerMap{
		conn:                       conn,
		connIDLen:                  connIDLen,
		handlers:                   make(map[string]packetHandlerEntry),
		resetTokens:                make(map[[16]byte]packetHandler),
		deleteRetiredSessionsAfter: protocol.RetiredConnectionIDDeleteTimeout,
		logger:                     logger,
	}
	go m.listen()
	return m
}

func (h *packetHandlerMap) Add(id protocol.ConnectionID, handler packetHandler) {
	h.mutex.Lock()
	h.handlers[string(id)] = packetHandlerEntry{handler: handler}
	h.mutex.Unlock()
}

func (h *packetHandlerMap) AddWithResetToken(id protocol.ConnectionID, handler packetHandler, token [16]byte) {
	h.mutex.Lock()
	h.handlers[string(id)] = packetHandlerEntry{handler: handler, resetToken: &token}
	h.resetTokens[token] = handler
	h.mutex.Unlock()
}

func (h *packetHandlerMap) Remove(id protocol.ConnectionID) {
	h.removeByConnectionIDAsString(string(id))
}

func (h *packetHandlerMap) removeByConnectionIDAsString(id string) {
	h.mutex.Lock()
	if handlerEntry, ok := h.handlers[id]; ok {
		if token := handlerEntry.resetToken; token != nil {
			delete(h.resetTokens, *token)
		}
		delete(h.handlers, id)
	}
	h.mutex.Unlock()
}

func (h *packetHandlerMap) Retire(id protocol.ConnectionID) {
	h.retireByConnectionIDAsString(string(id))
}

func (h *packetHandlerMap) retireByConnectionIDAsString(id string) {
	time.AfterFunc(h.deleteRetiredSessionsAfter, func() {
		h.removeByConnectionIDAsString(id)
	})
}

func (h *packetHandlerMap) SetServer(s unknownPacketHandler) {
	h.mutex.Lock()
	h.server = s
	h.mutex.Unlock()
}

func (h *packetHandlerMap) CloseServer() {
	h.mutex.Lock()
	h.server = nil
	var wg sync.WaitGroup
	for id, handlerEntry := range h.handlers {
		handler := handlerEntry.handler
		if handler.GetPerspective() == protocol.PerspectiveServer {
			wg.Add(1)
			go func(id string, handler packetHandler) {
				// session.Close() blocks until the CONNECTION_CLOSE has been sent and the run-loop has stopped
				_ = handler.Close()
				h.retireByConnectionIDAsString(id)
				wg.Done()
			}(id, handler)
		}
	}
	h.mutex.Unlock()
	wg.Wait()
}

func (h *packetHandlerMap) close(e error) error {
	h.mutex.Lock()
	if h.closed {
		h.mutex.Unlock()
		return nil
	}
	h.closed = true

	var wg sync.WaitGroup
	for _, handlerEntry := range h.handlers {
		wg.Add(1)
		go func(handlerEntry packetHandlerEntry) {
			handlerEntry.handler.destroy(e)
			wg.Done()
		}(handlerEntry)
	}

	if h.server != nil {
		h.server.closeWithError(e)
	}
	h.mutex.Unlock()
	wg.Wait()
	return getMultiplexer().RemoveConn(h.conn)
}

func (h *packetHandlerMap) listen() {
	for {
		buffer := getPacketBuffer()
		data := buffer.Slice
		// The packet size should not exceed protocol.MaxReceivePacketSize bytes
		// If it does, we only read a truncated packet, which will then end up undecryptable
		n, addr, err := h.conn.ReadFrom(data)
		if err != nil {
			h.close(err)
			return
		}
		h.handlePacket(addr, buffer, data[:n])
	}
}

func (h *packetHandlerMap) handlePacket(
	addr net.Addr,
	buffer *packetBuffer,
	data []byte,
) {
	packets, err := h.parsePacket(addr, buffer, data)
	if err != nil {
		h.logger.Debugf("error parsing packets from %s: %s", addr, err)
		// This is just the error from parsing the last packet.
		// We still need to process the packets that were successfully parsed before.
	}
	if len(packets) == 0 {
		buffer.Release()
		return
	}
	h.handleParsedPackets(packets)
}

func (h *packetHandlerMap) parsePacket(
	addr net.Addr,
	buffer *packetBuffer,
	data []byte,
) ([]*receivedPacket, error) {
	rcvTime := time.Now()
	packets := make([]*receivedPacket, 0, 1)

	var counter int
	var lastConnID protocol.ConnectionID
	for len(data) > 0 {
		hdr, err := wire.ParseHeader(bytes.NewReader(data), h.connIDLen)
		// drop the packet if we can't parse the header
		if err != nil {
			return packets, fmt.Errorf("error parsing header: %s", err)
		}
		if counter > 0 && !hdr.DestConnectionID.Equal(lastConnID) {
			return packets, fmt.Errorf("coalesced packet has different destination connection ID: %s, expected %s", hdr.DestConnectionID, lastConnID)
		}
		lastConnID = hdr.DestConnectionID

		var rest []byte
		if hdr.IsLongHeader {
			if protocol.ByteCount(len(data)) < hdr.ParsedLen()+hdr.Length {
				return packets, fmt.Errorf("packet length (%d bytes) is smaller than the expected length (%d bytes)", len(data)-int(hdr.ParsedLen()), hdr.Length)
			}
			packetLen := int(hdr.ParsedLen() + hdr.Length)
			rest = data[packetLen:]
			data = data[:packetLen]
		}

		if counter > 0 {
			buffer.Split()
		}
		counter++
		packets = append(packets, &receivedPacket{
			remoteAddr: addr,
			hdr:        hdr,
			rcvTime:    rcvTime,
			data:       data,
			buffer:     buffer,
		})

		// only log if this actually a coalesced packet
		if h.logger.Debug() && (counter > 1 || len(rest) > 0) {
			h.logger.Debugf("Parsed a coalesced packet. Part %d: %d bytes. Remaining: %d bytes.", counter, len(packets[counter-1].data), len(rest))
		}

		data = rest
	}
	return packets, nil
}

func (h *packetHandlerMap) handleParsedPackets(packets []*receivedPacket) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	// coalesced packets all have the same destination connection ID
	handlerEntry, handlerFound := h.handlers[string(packets[0].hdr.DestConnectionID)]

	for _, p := range packets {
		if handlerFound { // existing session
			handlerEntry.handler.handlePacket(p)
			continue
		}
		// No session found.
		// This might be a stateless reset.
		if !p.hdr.IsLongHeader {
			if len(p.data) >= protocol.MinStatelessResetSize {
				var token [16]byte
				copy(token[:], p.data[len(p.data)-16:])
				if sess, ok := h.resetTokens[token]; ok {
					sess.destroy(errors.New("received a stateless reset"))
					continue
				}
			}
			// TODO(#943): send a stateless reset
			h.logger.Debugf("received a short header packet with an unexpected connection ID %s", p.hdr.DestConnectionID)
			break // a short header packet is always the last in a coalesced packet
		}
		if h.server == nil { // no server set
			h.logger.Debugf("received a packet with an unexpected connection ID %s", p.hdr.DestConnectionID)
			continue
		}
		h.server.handlePacket(p)
	}
}
