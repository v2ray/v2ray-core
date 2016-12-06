package kcp

import (
	"crypto/tls"
	"net"
	"sync"
	"sync/atomic"
	"v2ray.com/core/common/alloc"
	"v2ray.com/core/common/dice"
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
	sync.Mutex
	net.Conn
	id    internal.ConnectionId
	input func([]byte)
	auth  internet.Authenticator
}

func (o *ClientConnection) Read([]byte) (int, error) {
	panic("KCP|ClientConnection: Read should not be called.")
}

func (o *ClientConnection) Id() internal.ConnectionId {
	return o.id
}

func (o *ClientConnection) Close() error {
	return o.Conn.Close()
}

func (o *ClientConnection) Reset(auth internet.Authenticator, inputCallback func([]byte)) {
	o.Lock()
	o.input = inputCallback
	o.auth = auth
	o.Unlock()
}

func (o *ClientConnection) Run() {
	payload := alloc.NewSmallBuffer()
	defer payload.Release()

	for {
		payload.Clear()
		_, err := payload.FillFrom(o.Conn)
		if err != nil {
			payload.Release()
			return
		}
		o.Lock()
		if o.input != nil && o.auth.Open(payload) {
			o.input(payload.Bytes())
		}
		o.Unlock()
		payload.Reset()
	}
}

func DialKCP(src v2net.Address, dest v2net.Destination, options internet.DialerOptions) (internet.Connection, error) {
	dest.Network = v2net.Network_UDP
	log.Info("KCP|Dialer: Dialing KCP to ", dest)

	id := internal.NewConnectionId(src, dest)
	conn := globalPool.Get(id)
	if conn == nil {
		rawConn, err := internet.DialToDest(src, dest)
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

	cpip, err := kcpSettings.GetAuthenticator()
	if err != nil {
		log.Error("KCP|Dialer: Failed to create authenticator: ", err)
		return nil, err
	}
	conv := uint16(atomic.AddUint32(&globalConv, 1))
	session := NewConnection(conv, conn.(*ClientConnection), globalPool, cpip, kcpSettings)

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
			tlsConn := tls.Client(conn, config)
			iConn = v2tls.NewConnection(tlsConn)
		}
	}

	return iConn, nil
}

func init() {
	internet.KCPDialer = DialKCP
}
