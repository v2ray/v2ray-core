package tcp

import (
	"context"

	"v2ray.com/core/app/log"
	"v2ray.com/core/common"
	"v2ray.com/core/common/net"
	"v2ray.com/core/transport/internet"
	"v2ray.com/core/transport/internet/tls"
)

func getTCPSettingsFromContext(ctx context.Context) *Config {
	rawTCPSettings := internet.TransportSettingsFromContext(ctx)
	if rawTCPSettings == nil {
		return nil
	}
	return rawTCPSettings.(*Config)
}

func Dial(ctx context.Context, dest net.Destination) (internet.Connection, error) {
	log.Trace(newError("dialing TCP to ", dest))
	src := internet.DialerSourceFromContext(ctx)

	conn, err := internet.DialSystem(ctx, src, dest)
	if err != nil {
		return nil, err
	}

	if config := tls.ConfigFromContext(ctx, tls.WithDestination(dest)); config != nil {
		conn = tls.Client(conn, config.GetTLSConfig())
	}

	tcpSettings := getTCPSettingsFromContext(ctx)
	if tcpSettings != nil && tcpSettings.HeaderSettings != nil {
		headerConfig, err := tcpSettings.HeaderSettings.GetInstance()
		if err != nil {
			return nil, newError("failed to get header settings").Base(err).AtError()
		}
		auth, err := internet.CreateConnectionAuthenticator(headerConfig)
		if err != nil {
			return nil, newError("failed to create header authenticator").Base(err).AtError()
		}
		conn = auth.Client(conn)
	}
	return internet.Connection(conn), nil
}

func init() {
	common.Must(internet.RegisterTransportDialer(internet.TransportProtocol_TCP, Dial))
}
