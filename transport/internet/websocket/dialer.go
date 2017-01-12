package websocket

import (
	"io/ioutil"
	"net"

	"github.com/gorilla/websocket"
	"v2ray.com/core/common"
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
	log.Info("WebSocket|Dailer: Creating connection to ", dest)
	if src == nil {
		src = v2net.AnyIP
	}
	networkSettings, err := options.Stream.GetEffectiveTransportSettings()
	if err != nil {
		return nil, err
	}
	wsSettings := networkSettings.(*Config)

	id := internal.NewConnectionID(src, dest)
	var conn *wsconn
	if dest.Network == v2net.Network_TCP && wsSettings.IsConnectionReuse() {
		connt := globalCache.Get(id)
		if connt != nil {
			conn = connt.(*wsconn)
		}
	}
	if conn == nil {
		var err error
		conn, err = wsDial(src, dest, options)
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

func wsDial(src v2net.Address, dest v2net.Destination, options internet.DialerOptions) (*wsconn, error) {
	networkSettings, err := options.Stream.GetEffectiveTransportSettings()
	if err != nil {
		return nil, err
	}
	wsSettings := networkSettings.(*Config)

	commonDial := func(network, addr string) (net.Conn, error) {
		return internet.DialSystem(src, dest)
	}

	dialer := websocket.Dialer{
		NetDial:         commonDial,
		ReadBufferSize:  65536,
		WriteBufferSize: 65536,
	}

	protocol := "ws"

	if options.Stream != nil && options.Stream.HasSecuritySettings() {
		protocol = "wss"
		securitySettings, err := options.Stream.GetEffectiveSecuritySettings()
		if err != nil {
			log.Error("WebSocket: Failed to create security settings: ", err)
			return nil, err
		}
		tlsConfig, ok := securitySettings.(*v2tls.Config)
		if ok {
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
	uri := protocol + "://" + host + "/" + wsSettings.Path

	conn, resp, err := dialer.Dial(uri, nil)
	if err != nil {
		if resp != nil {
			reason, reasonerr := ioutil.ReadAll(resp.Body)
			log.Info(string(reason), reasonerr)
		}
		return nil, err
	}
	return func() internet.Connection {
		connv2ray := &wsconn{
			wsc:         conn,
			connClosing: false,
			config:      wsSettings,
		}
		connv2ray.setup()
		return connv2ray
	}().(*wsconn), nil
}
