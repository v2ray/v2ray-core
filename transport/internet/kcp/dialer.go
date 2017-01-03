package kcp

import (
	"crypto/tls"
	"net"
	"sync"
	"sync/atomic"

	"crypto/cipher"

	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/dice"
	"v2ray.com/core/common/errors"
	"v2ray.com/core/common/log"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/transport/internet"
	"v2ray.com/core/transport/internet/internal"
	v2tls "v2ray.com/core/transport/internet/tls"
)

var (
	globalConv = uint32(dice.Roll(65536))
	globalPool = internal.NewConnectionPool()
)

type ClientConnection struct {
	sync.RWMutex
	net.Conn
	id     internal.ConnectionID
	input  func([]Segment)
	reader PacketReader
	writer PacketWriter
}

func (o *ClientConnection) Overhead() int {
	o.RLock()
	defer o.RUnlock()
	if o.writer == nil {
		return 0
	}
	return o.writer.Overhead()
}

func (o *ClientConnection) Write(b []byte) (int, error) {
	o.RLock()
	defer o.RUnlock()

	if o.writer == nil {
		return len(b), nil
	}

	return o.writer.Write(b)
}

func (o *ClientConnection) Read([]byte) (int, error) {
	panic("KCP|ClientConnection: Read should not be called.")
}

func (o *ClientConnection) Id() internal.ConnectionID {
	return o.id
}

func (o *ClientConnection) Close() error {
	return o.Conn.Close()
}

func (o *ClientConnection) Reset(inputCallback func([]Segment)) {
	o.Lock()
	o.input = inputCallback
	o.Unlock()
}

func (o *ClientConnection) ResetSecurity(header internet.PacketHeader, security cipher.AEAD) {
	o.Lock()
	if o.reader == nil {
		o.reader = new(KCPPacketReader)
	}
	o.reader.(*KCPPacketReader).Header = header
	o.reader.(*KCPPacketReader).Security = security
	if o.writer == nil {
		o.writer = new(KCPPacketWriter)
	}
	o.writer.(*KCPPacketWriter).Header = header
	o.writer.(*KCPPacketWriter).Security = security
	o.writer.(*KCPPacketWriter).Writer = o.Conn

	o.Unlock()
}

func (o *ClientConnection) Run() {
	payload := buf.NewSmall()
	defer payload.Release()

	for {
		err := payload.Reset(buf.ReadFrom(o.Conn))
		if err != nil {
			payload.Release()
			return
		}
		o.RLock()
		if o.input != nil {
			segments := o.reader.Read(payload.Bytes())
			if len(segments) > 0 {
				o.input(segments)
			}
		}
		o.RUnlock()
	}
}

func DialKCP(src v2net.Address, dest v2net.Destination, options internet.DialerOptions) (internet.Connection, error) {
	dest.Network = v2net.Network_UDP
	log.Info("KCP|Dialer: Dialing KCP to ", dest)

	id := internal.NewConnectionID(src, dest)
	conn := globalPool.Get(id)
	if conn == nil {
		rawConn, err := internet.DialSystem(src, dest)
		if err != nil {
			log.Error("KCP|Dialer: Failed to dial to dest: ", err)
			return nil, err
		}
		c := &ClientConnection{
			Conn: rawConn,
			id:   id,
		}
		go c.Run()
		conn = c
	}

	networkSettings, err := options.Stream.GetEffectiveNetworkSettings()
	if err != nil {
		log.Error("KCP|Dialer: Failed to get KCP settings: ", err)
		return nil, err
	}
	kcpSettings := networkSettings.(*Config)

	clientConn := conn.(*ClientConnection)
	header, err := kcpSettings.GetPackerHeader()
	if err != nil {
		return nil, errors.Base(err).Message("KCP|Dialer: Failed to create packet header.")
	}
	security, err := kcpSettings.GetSecurity()
	if err != nil {
		return nil, errors.Base(err).Message("KCP|Dialer: Failed to create security.")
	}
	clientConn.ResetSecurity(header, security)
	conv := uint16(atomic.AddUint32(&globalConv, 1))
	session := NewConnection(conv, clientConn, globalPool, kcpSettings)

	var iConn internet.Connection
	iConn = session

	if options.Stream != nil && options.Stream.HasSecuritySettings() {
		securitySettings, err := options.Stream.GetEffectiveSecuritySettings()
		if err != nil {
			log.Error("KCP|Dialer: Failed to get security settings: ", err)
			return nil, err
		}
		switch securitySettings := securitySettings.(type) {
		case *v2tls.Config:
			config := securitySettings.GetTLSConfig()
			if dest.Address.Family().IsDomain() {
				config.ServerName = dest.Address.Domain()
			}
			tlsConn := tls.Client(iConn, config)
			iConn = v2tls.NewConnection(tlsConn)
		}
	}

	return iConn, nil
}

func init() {
	common.Must(internet.RegisterNetworkDialer(v2net.Network_KCP, DialKCP))
}
