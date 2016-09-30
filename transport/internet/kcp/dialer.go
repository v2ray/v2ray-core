package kcp

import (
	"crypto/tls"
	"net"
	"sync/atomic"

	"v2ray.com/core/common/dice"
	"v2ray.com/core/common/log"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/transport/internet"
	v2tls "v2ray.com/core/transport/internet/tls"
)

var (
	globalConv = uint32(dice.Roll(65536))
)

func DialKCP(src v2net.Address, dest v2net.Destination, options internet.DialerOptions) (internet.Connection, error) {
	dest.Network = v2net.Network_UDP
	log.Info("KCP|Dialer: Dialing KCP to ", dest)
	conn, err := internet.DialToDest(src, dest)
	if err != nil {
		log.Error("KCP|Dialer: Failed to dial to dest: ", err)
		return nil, err
	}

	cpip, err := effectiveConfig.GetAuthenticator()
	if err != nil {
		log.Error("KCP|Dialer: Failed to create authenticator: ", err)
		return nil, err
	}
	conv := uint16(atomic.AddUint32(&globalConv, 1))
	session := NewConnection(conv, conn, conn.LocalAddr().(*net.UDPAddr), conn.RemoteAddr().(*net.UDPAddr), cpip)
	session.FetchInputFrom(conn)

	var iConn internet.Connection
	iConn = session

	if options.Stream != nil && options.Stream.Security == internet.StreamSecurityTypeTLS {
		config := options.Stream.TLSSettings.GetTLSConfig()
		if dest.Address.Family().IsDomain() {
			config.ServerName = dest.Address.Domain()
		}
		tlsConn := tls.Client(conn, config)
		iConn = v2tls.NewConnection(tlsConn)
	}

	return iConn, nil
}

func init() {
	internet.KCPDialer = DialKCP
}
