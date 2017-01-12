package tcp

import (
	"crypto/tls"
	"net"

	"v2ray.com/core/common"
	"v2ray.com/core/common/errors"
	"v2ray.com/core/common/log"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/transport/internet"
	"v2ray.com/core/transport/internet/internal"
	v2tls "v2ray.com/core/transport/internet/tls"
)

var (
	globalCache = internal.NewConnectionPool()
)

func Dial(src v2net.Address, dest v2net.Destination, options internet.DialerOptions) (internet.Connection, error) {
	log.Info("Internet|TCP: Dailing TCP to ", dest)
	if src == nil {
		src = v2net.AnyIP
	}
	networkSettings, err := options.Stream.GetEffectiveTransportSettings()
	if err != nil {
		return nil, err
	}
	tcpSettings := networkSettings.(*Config)

	id := internal.NewConnectionID(src, dest)
	var conn net.Conn
	if dest.Network == v2net.Network_TCP && tcpSettings.IsConnectionReuse() {
		conn = globalCache.Get(id)
	}
	if conn == nil {
		var err error
		conn, err = internet.DialSystem(src, dest)
		if err != nil {
			return nil, err
		}
		if options.Stream != nil && options.Stream.HasSecuritySettings() {
			securitySettings, err := options.Stream.GetEffectiveSecuritySettings()
			if err != nil {
				log.Error("TCP: Failed to get security settings: ", err)
				return nil, err
			}
			tlsConfig, ok := securitySettings.(*v2tls.Config)
			if ok {
				config := tlsConfig.GetTLSConfig()
				if dest.Address.Family().IsDomain() {
					config.ServerName = dest.Address.Domain()
				}
				conn = tls.Client(conn, config)
			}
		}
		if tcpSettings.HeaderSettings != nil {
			headerConfig, err := tcpSettings.HeaderSettings.GetInstance()
			if err != nil {
				return nil, errors.Base(err).Message("Interent|TCP: Failed to get header settings.")
			}
			auth, err := internet.CreateConnectionAuthenticator(tcpSettings.HeaderSettings.Type, headerConfig)
			if err != nil {
				return nil, errors.Base(err).Message("Internet|TCP: Failed to create header authenticator.")
			}
			conn = auth.Client(conn)
		}
	}
	return internal.NewConnection(id, conn, globalCache, internal.ReuseConnection(tcpSettings.IsConnectionReuse())), nil
}

func init() {
	common.Must(internet.RegisterTransportDialer(internet.TransportProtocol_TCP, Dial))
}
