package hub

import (
	"errors"
	"net"
	"time"

	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/proxy"
	"github.com/v2ray/v2ray-core/transport"
)

var (
	ErrorInvalidHost = errors.New("Invalid Host.")

	globalCache = NewConnectionCache()
)

func Dial(src v2net.Address, dest v2net.Destination) (*Connection, error) {
	if src == nil {
		src = v2net.AnyIP
	}
	id := src.String() + "-" + dest.NetAddr()
	var conn net.Conn
	if dest.IsTCP() && transport.IsConnectionReusable() {
		conn = globalCache.Get(id)
	}
	if conn == nil {
		var err error
		conn, err = DialWithoutCache(src, dest)
		if err != nil {
			return nil, err
		}
	}
	return &Connection{
		dest:     id,
		conn:     conn,
		listener: globalCache,
	}, nil
}

func DialWithoutCache(src v2net.Address, dest v2net.Destination) (net.Conn, error) {
	dialer := &net.Dialer{
		Timeout:   time.Second * 60,
		DualStack: true,
	}

	if src != nil && src != v2net.AnyIP {
		var addr net.Addr
		if dest.IsTCP() {
			addr = &net.TCPAddr{
				IP:   src.IP(),
				Port: 0,
			}
		} else {
			addr = &net.UDPAddr{
				IP:   src.IP(),
				Port: 0,
			}
		}
		dialer.LocalAddr = addr
	}

	return dialer.Dial(dest.Network().String(), dest.NetAddr())
}

func Dial3(src v2net.Address, dest v2net.Destination, proxyMeta *proxy.OutboundHandlerMeta) (*Connection, error) {
	if proxyMeta.KcpSupported && transport.IsKcpEnabled() {
		DialKCP3(src, dest, proxyMeta)
	}
	return Dial(src, dest)
}
func DialWithoutCache3(src v2net.Address, dest v2net.Destination, proxyMeta *proxy.OutboundHandlerMeta) (net.Conn, error) {
	if proxyMeta.KcpSupported && transport.IsKcpEnabled() {
	}
	return DialWithoutCache(src, dest)
}

func DialKCP3(src v2net.Address, dest v2net.Destination, proxyMeta *proxy.OutboundHandlerMeta) (*Connection, error) {
	if src == nil {
		src = v2net.AnyIP
	}
	id := src.String() + "-" + dest.NetAddr()
	var conn net.Conn
	if dest.IsTCP() && transport.IsConnectionReusable() {
		conn = globalCache.Get(id)
	}
	if conn == nil {
		var err error
		conn, err = DialWithoutCache3(src, dest, proxyMeta)
		if err != nil {
			return nil, err
		}
	}
	return &Connection{
		dest:     id,
		conn:     conn,
		listener: globalCache,
	}, nil
}

/*DialKCPWithoutCache Dial KCP connection
This Dialer will ignore src this is a restriction
due to github.com/xtaci/kcp-go.DialWithOptions
*/
func DialKCPWithoutCache(src v2net.Address, dest v2net.Destination) (net.Conn, error) {
	return DialKCP(dest)
}
