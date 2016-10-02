package tcp

import (
	"net"

	"crypto/tls"
	"v2ray.com/core/common/log"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/transport/internet"
	v2tls "v2ray.com/core/transport/internet/tls"
)

var (
	globalCache = NewConnectionCache()
)

func Dial(src v2net.Address, dest v2net.Destination, options internet.DialerOptions) (internet.Connection, error) {
	log.Info("Dailing TCP to ", dest)
	if src == nil {
		src = v2net.AnyIP
	}
	networkSettings, err := options.Stream.GetEffectiveNetworkSettings()
	if err != nil {
		return nil, err
	}
	tcpSettings := networkSettings.(*Config)

	id := src.String() + "-" + dest.NetAddr()
	var conn net.Conn
	if dest.Network == v2net.Network_TCP && tcpSettings.ConnectionReuse {
		conn = globalCache.Get(id)
	}
	if conn == nil {
		var err error
		conn, err = internet.DialToDest(src, dest)
		if err != nil {
			return nil, err
		}
	}
	if options.Stream != nil && options.Stream.SecurityType == internet.SecurityType_TLS {
		securitySettings, err := options.Stream.GetEffectiveSecuritySettings()
		if err != nil {
			log.Error("TCP: Failed to apply TLS config: ", err)
			return nil, err
		}
		config := securitySettings.(*v2tls.Config).GetTLSConfig()
		if dest.Address.Family().IsDomain() {
			config.ServerName = dest.Address.Domain()
		}
		conn = tls.Client(conn, config)
	}
	return NewConnection(id, conn, globalCache, tcpSettings), nil
}

func DialRaw(src v2net.Address, dest v2net.Destination, options internet.DialerOptions) (internet.Connection, error) {
	log.Info("Dailing Raw TCP to ", dest)
	conn, err := internet.DialToDest(src, dest)
	if err != nil {
		return nil, err
	}
	// TODO: handle dialer options
	return &RawConnection{
		TCPConn: *conn.(*net.TCPConn),
	}, nil
}

func init() {
	internet.TCPDialer = Dial
	internet.RawTCPDialer = DialRaw
}
