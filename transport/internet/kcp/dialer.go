package kcp

import (
	"context"
	"crypto/cipher"
	"crypto/tls"
	"net"
	"sync"
	"sync/atomic"

	"v2ray.com/core/app/log"
	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/dice"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/transport/internet"
	v2tls "v2ray.com/core/transport/internet/tls"
)

var (
	globalConv = uint32(dice.RollUint16())
)

type ClientConnection struct {
	sync.RWMutex
	net.Conn
	input  func([]Segment)
	reader PacketReader
	writer PacketWriter
}

func (c *ClientConnection) Overhead() int {
	c.RLock()
	defer c.RUnlock()
	if c.writer == nil {
		return 0
	}
	return c.writer.Overhead()
}

// Write implements io.Writer.
func (c *ClientConnection) Write(b []byte) (int, error) {
	c.RLock()
	defer c.RUnlock()

	if c.writer == nil {
		return len(b), nil
	}

	return c.writer.Write(b)
}

func (*ClientConnection) Read([]byte) (int, error) {
	panic("KCP|ClientConnection: Read should not be called.")
}

func (c *ClientConnection) Close() error {
	return c.Conn.Close()
}

func (c *ClientConnection) Reset(inputCallback func([]Segment)) {
	c.Lock()
	c.input = inputCallback
	c.Unlock()
}

func (c *ClientConnection) ResetSecurity(header internet.PacketHeader, security cipher.AEAD) {
	c.Lock()
	if c.reader == nil {
		c.reader = new(KCPPacketReader)
	}
	c.reader.(*KCPPacketReader).Header = header
	c.reader.(*KCPPacketReader).Security = security
	if c.writer == nil {
		c.writer = new(KCPPacketWriter)
	}
	c.writer.(*KCPPacketWriter).Header = header
	c.writer.(*KCPPacketWriter).Security = security
	c.writer.(*KCPPacketWriter).Writer = c.Conn

	c.Unlock()
}

func (c *ClientConnection) Run() {
	payload := buf.New()
	defer payload.Release()

	for {
		err := payload.Reset(buf.ReadFrom(c.Conn))
		if err != nil {
			payload.Release()
			return
		}
		c.RLock()
		if c.input != nil {
			segments := c.reader.Read(payload.Bytes())
			if len(segments) > 0 {
				c.input(segments)
			}
		}
		c.RUnlock()
	}
}

func DialKCP(ctx context.Context, dest v2net.Destination) (internet.Connection, error) {
	dest.Network = v2net.Network_UDP
	log.Trace(newError("dialing mKCP to ", dest))

	src := internet.DialerSourceFromContext(ctx)
	rawConn, err := internet.DialSystem(ctx, src, dest)
	if err != nil {
		return nil, newError("failed to dial to dest: ", err).AtWarning().Base(err)
	}
	conn := &ClientConnection{
		Conn: rawConn,
	}
	go conn.Run()

	kcpSettings := internet.TransportSettingsFromContext(ctx).(*Config)

	header, err := kcpSettings.GetPackerHeader()
	if err != nil {
		return nil, newError("failed to create packet header").Base(err)
	}
	security, err := kcpSettings.GetSecurity()
	if err != nil {
		return nil, newError("failed to create security").Base(err)
	}
	conn.ResetSecurity(header, security)
	conv := uint16(atomic.AddUint32(&globalConv, 1))
	session := NewConnection(conv, conn, kcpSettings)

	var iConn internet.Connection
	iConn = session

	if securitySettings := internet.SecuritySettingsFromContext(ctx); securitySettings != nil {
		switch securitySettings := securitySettings.(type) {
		case *v2tls.Config:
			config := securitySettings.GetTLSConfig()
			if dest.Address.Family().IsDomain() {
				config.ServerName = dest.Address.Domain()
			}
			tlsConn := tls.Client(iConn, config)
			iConn = tlsConn
		}
	}

	return iConn, nil
}

func init() {
	common.Must(internet.RegisterTransportDialer(internet.TransportProtocol_MKCP, DialKCP))
}
