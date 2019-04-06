// +build !confonly

package quic

import (
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"time"

	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/net"
	quic "v2ray.com/core/external/github.com/lucas-clemente/quic-go"
	"v2ray.com/core/transport/internet"
)

type sysConn struct {
	conn   net.PacketConn
	header internet.PacketHeader
	auth   cipher.AEAD
}

func wrapSysConn(rawConn net.PacketConn, config *Config) (*sysConn, error) {
	header, err := getHeader(config)
	if err != nil {
		return nil, err
	}
	auth, err := getAuth(config)
	if err != nil {
		return nil, err
	}
	return &sysConn{
		conn:   rawConn,
		header: header,
		auth:   auth,
	}, nil
}

var errInvalidPacket = errors.New("invalid packet")

func (c *sysConn) readFromInternal(p []byte) (int, net.Addr, error) {
	buffer := getBuffer()
	defer putBuffer(buffer)

	nBytes, addr, err := c.conn.ReadFrom(buffer)
	if err != nil {
		return 0, nil, err
	}

	payload := buffer[:nBytes]
	if c.header != nil {
		if len(payload) <= int(c.header.Size()) {
			return 0, nil, errInvalidPacket
		}
		payload = payload[c.header.Size():]
	}

	if c.auth == nil {
		n := copy(p, payload)
		return n, addr, nil
	}

	if len(payload) <= c.auth.NonceSize() {
		return 0, nil, errInvalidPacket
	}

	nonce := payload[:c.auth.NonceSize()]
	payload = payload[c.auth.NonceSize():]

	p, err = c.auth.Open(p[:0], nonce, payload, nil)
	if err != nil {
		return 0, nil, errInvalidPacket
	}

	return len(p), addr, nil
}

func (c *sysConn) ReadFrom(p []byte) (int, net.Addr, error) {
	if c.header == nil && c.auth == nil {
		return c.conn.ReadFrom(p)
	}

	for {
		n, addr, err := c.readFromInternal(p)
		if err != nil && err != errInvalidPacket {
			return 0, nil, err
		}
		if err == nil {
			return n, addr, nil
		}
	}
}

func (c *sysConn) WriteTo(p []byte, addr net.Addr) (int, error) {
	if c.header == nil && c.auth == nil {
		return c.conn.WriteTo(p, addr)
	}

	buffer := getBuffer()
	defer putBuffer(buffer)

	payload := buffer
	n := 0
	if c.header != nil {
		c.header.Serialize(payload)
		n = int(c.header.Size())
	}

	if c.auth == nil {
		nBytes := copy(payload[n:], p)
		n += nBytes
	} else {
		nounce := payload[n : n+c.auth.NonceSize()]
		common.Must2(rand.Read(nounce))
		n += c.auth.NonceSize()
		pp := c.auth.Seal(payload[:n], nounce, p, nil)
		n = len(pp)
	}

	return c.conn.WriteTo(payload[:n], addr)
}

func (c *sysConn) Close() error {
	return c.conn.Close()
}

func (c *sysConn) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

func (c *sysConn) SetDeadline(t time.Time) error {
	return c.conn.SetDeadline(t)
}

func (c *sysConn) SetReadDeadline(t time.Time) error {
	return c.conn.SetReadDeadline(t)
}

func (c *sysConn) SetWriteDeadline(t time.Time) error {
	return c.conn.SetWriteDeadline(t)
}

type interConn struct {
	stream quic.Stream
	local  net.Addr
	remote net.Addr
}

func (c *interConn) Read(b []byte) (int, error) {
	return c.stream.Read(b)
}

func (c *interConn) WriteMultiBuffer(mb buf.MultiBuffer) error {
	mb = buf.Compact(mb)
	mb, err := buf.WriteMultiBuffer(c, mb)
	buf.ReleaseMulti(mb)
	return err
}

func (c *interConn) Write(b []byte) (int, error) {
	return c.stream.Write(b)
}

func (c *interConn) Close() error {
	return c.stream.Close()
}

func (c *interConn) LocalAddr() net.Addr {
	return c.local
}

func (c *interConn) RemoteAddr() net.Addr {
	return c.remote
}

func (c *interConn) SetDeadline(t time.Time) error {
	return c.stream.SetDeadline(t)
}

func (c *interConn) SetReadDeadline(t time.Time) error {
	return c.stream.SetReadDeadline(t)
}

func (c *interConn) SetWriteDeadline(t time.Time) error {
	return c.stream.SetWriteDeadline(t)
}
