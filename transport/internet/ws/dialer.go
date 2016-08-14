package ws

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net"

	"github.com/gorilla/websocket"
	"github.com/v2ray/v2ray-core/common/log"
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/transport/internet"
)

var (
	globalCache = NewConnectionCache()
)

func Dial(src v2net.Address, dest v2net.Destination) (internet.Connection, error) {
	log.Info("Dailing WS to ", dest)
	if src == nil {
		src = v2net.AnyIP
	}
	id := src.String() + "-" + dest.NetAddr()
	var conn *wsconn
	if dest.IsTCP() && effectiveConfig.ConnectionReuse {
		connt := globalCache.Get(id)
		if connt != nil {
			conn = connt.(*wsconn)
		}
	}
	if conn == nil {
		var err error
		conn, err = wsDial(src, dest)
		if err != nil {
			log.Warning("WS Dial failed:" + err.Error())
			return nil, err
		}
	}
	return NewConnection(id, conn, globalCache), nil
}

func init() {
	internet.WSDialer = Dial
}

func wsDial(src v2net.Address, dest v2net.Destination) (*wsconn, error) {
	commonDial := func(network, addr string) (net.Conn, error) {
		return internet.DialToDest(src, dest)
	}

	tlsconf := &tls.Config{ServerName: dest.Address().Domain()}

	dialer := websocket.Dialer{NetDial: commonDial, ReadBufferSize: 65536, WriteBufferSize: 65536, TLSClientConfig: tlsconf}

	effpto := func(dst v2net.Destination) string {

		if effectiveConfig.Pto != "" {
			return effectiveConfig.Pto
		}

		switch dst.Port().Value() {
		/*
				Since the value is not given explicitly,
				We are guessing it now.

				HTTP Port:
						80
				    8080
				    8880
				    2052
			      2082
			      2086
			      2095

				HTTPS Port:
						443
				    2053
				    2083
			      2087
			      2096
			      8443

				if the port you are using is not well-known,
				specify it to avoid this process.

				We will	return "CRASH"turn "unknown" if we can't guess it, cause Dial to fail.
		*/
		case 80:
		case 8080:
		case 8880:
		case 2052:
		case 2082:
		case 2086:
		case 2095:
			return "ws"
		case 443:
		case 2053:
		case 2083:
		case 2087:
		case 2096:
		case 8443:
			return "wss"
		default:
			return "unknown"
		}
		panic("Runtime unstable. Please report this bug to developers.")
	}(dest)

	uri := func(dst v2net.Destination, pto string, path string) string {
		return fmt.Sprintf("%v://%v:%v/%v", pto, dst.NetAddr(), dst.Port(), path)
	}(dest, effpto, effectiveConfig.Path)

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
