package tcp

import (
	"context"
	"crypto/tls"
	"net"

	"v2ray.com/core/app/log"
	"v2ray.com/core/common"
	"v2ray.com/core/common/errors"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/transport/internet"
	"v2ray.com/core/transport/internet/internal"
	v2tls "v2ray.com/core/transport/internet/tls"
)

var (
	globalCache = internal.NewConnectionPool()
)

func Dial(ctx context.Context, dest v2net.Destination) (internet.Connection, error) {
	log.Trace(errors.New("Internet|TCP: Dailing TCP to ", dest))
	src := internet.DialerSourceFromContext(ctx)

	tcpSettings := internet.TransportSettingsFromContext(ctx).(*Config)

	id := internal.NewConnectionID(src, dest)
	var conn net.Conn
	if dest.Network == v2net.Network_TCP && tcpSettings.IsConnectionReuse() {
		conn = globalCache.Get(id)
	}
	if conn == nil {
		var err error
		conn, err = internet.DialSystem(ctx, src, dest)
		if err != nil {
			return nil, err
		}
		if securitySettings := internet.SecuritySettingsFromContext(ctx); securitySettings != nil {
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
				return nil, errors.New("Internet|TCP: Failed to get header settings.").Base(err)
			}
			auth, err := internet.CreateConnectionAuthenticator(headerConfig)
			if err != nil {
				return nil, errors.New("Internet|TCP: Failed to create header authenticator.").Base(err)
			}
			conn = auth.Client(conn)
		}
	}
	return internal.NewConnection(id, conn, globalCache, internal.ReuseConnection(tcpSettings.IsConnectionReuse())), nil
}

func init() {
	common.Must(internet.RegisterTransportDialer(internet.TransportProtocol_TCP, Dial))
}
