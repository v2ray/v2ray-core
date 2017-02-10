package websocket

import (
	"context"
	"net"

	"github.com/gorilla/websocket"
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
	log.Info("WebSocket|Dialer: Creating connection to ", dest)
	src := internet.DialerSourceFromContext(ctx)
	wsSettings := internet.TransportSettingsFromContext(ctx).(*Config)

	id := internal.NewConnectionID(src, dest)
	var conn net.Conn
	if dest.Network == v2net.Network_TCP && wsSettings.IsConnectionReuse() {
		conn = globalCache.Get(id)
	}
	if conn == nil {
		var err error
		conn, err = wsDial(ctx, dest)
		if err != nil {
			log.Warning("WebSocket|Dialer: Dial failed: ", err)
			return nil, err
		}
	}
	return internal.NewConnection(id, conn, globalCache, internal.ReuseConnection(wsSettings.IsConnectionReuse())), nil
}

func init() {
	common.Must(internet.RegisterTransportDialer(internet.TransportProtocol_WebSocket, Dial))
}

func wsDial(ctx context.Context, dest v2net.Destination) (net.Conn, error) {
	src := internet.DialerSourceFromContext(ctx)
	wsSettings := internet.TransportSettingsFromContext(ctx).(*Config)

	commonDial := func(network, addr string) (net.Conn, error) {
		return internet.DialSystem(src, dest)
	}

	dialer := websocket.Dialer{
		NetDial:         commonDial,
		ReadBufferSize:  32 * 1024,
		WriteBufferSize: 32 * 1024,
	}

	protocol := "ws"

	if securitySettings := internet.SecuritySettingsFromContext(ctx); securitySettings != nil {
		tlsConfig, ok := securitySettings.(*v2tls.Config)
		if ok {
			protocol = "wss"
			dialer.TLSClientConfig = tlsConfig.GetTLSConfig()
			if dest.Address.Family().IsDomain() {
				dialer.TLSClientConfig.ServerName = dest.Address.Domain()
			}
		}
	}

	host := dest.NetAddr()
	if (protocol == "ws" && dest.Port == 80) || (protocol == "wss" && dest.Port == 443) {
		host = dest.Address.String()
	}
	uri := protocol + "://" + host + wsSettings.GetNormailzedPath()

	conn, resp, err := dialer.Dial(uri, nil)
	if err != nil {
		var reason string
		if resp != nil {
			reason = resp.Status
		}
		return nil, errors.Base(err).Message("WebSocket|Dialer: Failed to dial to (", uri, "): ", reason)
	}

	return &wsconn{
		wsc: conn,
	}, nil
}
