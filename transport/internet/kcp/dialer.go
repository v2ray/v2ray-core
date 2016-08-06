package kcp

import (
	"net"
	"sync/atomic"

	"github.com/v2ray/v2ray-core/common/dice"
	"github.com/v2ray/v2ray-core/common/log"
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/transport/internet"
)

var (
	globalConv = uint32(dice.Roll(65536))
)

func DialKCP(src v2net.Address, dest v2net.Destination) (internet.Connection, error) {
	udpDest := v2net.UDPDestination(dest.Address(), dest.Port())
	log.Info("KCP|Dialer: Dialing KCP to ", udpDest)
	conn, err := internet.DialToDest(src, udpDest)
	if err != nil {
		log.Error("KCP|Dialer: Failed to dial to dest: ", err)
		return nil, err
	}

	cpip, err := effectiveConfig.GetAuthenticator()
	if err != nil {
		log.Error("KCP|Dialer: Failed to create authenticator: ", err)
		return nil, err
	}
	conv := uint16(atomic.AddUint32(&globalConv, 1))
	session := NewConnection(conv, conn, conn.LocalAddr().(*net.UDPAddr), conn.RemoteAddr().(*net.UDPAddr), cpip)
	session.FetchInputFrom(conn)

	return session, nil
}

func init() {
	internet.KCPDialer = DialKCP
}
