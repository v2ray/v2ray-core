package kcp

import (
	"math/rand"
	"net"

	"github.com/v2ray/v2ray-core/common/log"
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/transport/internet"
)

func DialKCP(src v2net.Address, dest v2net.Destination) (internet.Connection, error) {
	udpDest := v2net.UDPDestination(dest.Address(), dest.Port())
	log.Info("Dialling KCP to ", udpDest)
	conn, err := internet.DialToDest(src, udpDest)
	if err != nil {
		return nil, err
	}

	cpip := NewSimpleAuthenticator()
	session := NewConnection(uint16(rand.Uint32()), conn, conn.LocalAddr().(*net.UDPAddr), conn.RemoteAddr().(*net.UDPAddr), cpip)
	session.FetchInputFrom(conn)

	return session, nil
}

func init() {
	internet.KCPDialer = DialKCP
}
