package ws

import (
	"fmt"
	"io/ioutil"
	"net"

	"github.com/gorilla/websocket"
	"v2ray.com/core/common/log"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/transport/internet"
)

var (
	globalCache = NewConnectionCache()
)

func Dial(src v2net.Address, dest v2net.Destination, options internet.DialerOptions) (internet.Connection, error) {
	log.Info("WebSocket|Dailer: Creating connection to ", dest)
	if src == nil {
		src = v2net.AnyIP
	}
	id := src.String() + "-" + dest.NetAddr()
	var conn *wsconn
	if dest.Network == v2net.Network_TCP && effectiveConfig.ConnectionReuse {
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
	return NewConnection(id, conn, globalCache), nil
}

func init() {
	internet.WSDialer = Dial
}

func wsDial(src v2net.Address, dest v2net.Destination, options internet.DialerOptions) (*wsconn, error) {
	commonDial := func(network, addr string) (net.Conn, error) {
		return internet.DialToDest(src, dest)
	}

	dialer := websocket.Dialer{
		NetDial:         commonDial,
		ReadBufferSize:  65536,
		WriteBufferSize: 65536,
	}

	protocol := "ws"

	if options.Stream != nil && options.Stream.Security == internet.StreamSecurityTypeTLS {
		protocol = "wss"
		dialer.TLSClientConfig = options.Stream.TLSSettings.GetTLSConfig()
		if dest.Address.Family().IsDomain() {
			dialer.TLSClientConfig.ServerName = dest.Address.Domain()
		}
	}

	uri := func(dst v2net.Destination, pto string, path string) string {
		return fmt.Sprintf("%v://%v/%v", pto, dst.NetAddr(), path)
	}(dest, protocol, effectiveConfig.Path)

	conn, resp, err := dialer.Dial(uri, nil)
	if err != nil {
		if resp != nil {
			reason, reasonerr := ioutil.ReadAll(resp.Body)
			log.Info(string(reason), reasonerr)
		}
		return nil, err
	}
	return func() internet.Connection {
		connv2ray := &wsconn{wsc: conn, connClosing: false}
		connv2ray.setup()
		return connv2ray
	}().(*wsconn), nil
}
