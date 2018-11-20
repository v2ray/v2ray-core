package quic

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Client Multiplexer", func() {
	It("adds a new packet conn ", func() {
		conn := newMockPacketConn()
		_, err := getMultiplexer().AddConn(conn, 8)
		Expect(err).ToNot(HaveOccurred())
	})

	It("errors when adding an existing conn with a different connection ID length", func() {
		conn := newMockPacketConn()
		_, err := getMultiplexer().AddConn(conn, 5)
		Expect(err).ToNot(HaveOccurred())
		_, err = getMultiplexer().AddConn(conn, 6)
		Expect(err).To(MatchError("cannot use 6 byte connection IDs on a connection that is already using 5 byte connction IDs"))
	})

})
