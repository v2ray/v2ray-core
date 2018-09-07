package tcp

import (
	"context"

	"v2ray.com/core/common"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/session"
	"v2ray.com/core/transport/internet"
	"v2ray.com/core/transport/internet/tls"
)

func getTCPSettingsFromContext(ctx context.Context) *Config {
	rawTCPSettings := internet.StreamSettingsFromContext(ctx)
	if rawTCPSettings == nil {
		return nil
	}
	return rawTCPSettings.ProtocolSettings.(*Config)
}

// Dial dials a new TCP connection to the given destination.
func Dial(ctx context.Context, dest net.Destination) (internet.Connection, error) {
	newError("dialing TCP to ", dest).WriteToLog(session.ExportIDToError(ctx))
	src := internet.DialerSourceFromContext(ctx)

	conn, err := internet.DialSystem(ctx, src, dest)
	if err != nil {
		return nil, err
	}

	if config := tls.ConfigFromContext(ctx); config != nil {
		conn = tls.Client(conn, config.GetTLSConfig(tls.WithDestination(dest), tls.WithNextProto("h2")))
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
	common.Must(internet.RegisterTransportDialer(protocolName, Dial))
}
