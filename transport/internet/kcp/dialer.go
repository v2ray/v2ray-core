package kcp

import (
	"errors"
	"math/rand"
	"net"

	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/transport/internet"
)

var (
	ErrUnknownDestination = errors.New("Destination IP can't be resolved.")
)

func DialKCP(src v2net.Address, dest v2net.Destination) (internet.Connection, error) {
	var ip net.IP
	if dest.Address().IsDomain() {
		ips, err := net.LookupIP(dest.Address().Domain())
		if err != nil {
			return nil, err
		}
		if len(ips) == 0 {
			return nil, ErrUnknownDestination
		}
		ip = ips[0]
	} else {
		ip = dest.Address().IP()
	}
	udpAddr := &net.UDPAddr{
		IP:   ip,
		Port: int(dest.Port()),
	}

	udpConn, err := net.ListenUDP("udp", &net.UDPAddr{})
	if err != nil {
		return nil, err
	}

	cpip, _ := NewNoneBlockCrypt(nil)
	session := newUDPSession(rand.Uint32(), nil, udpConn, udpAddr, cpip)
	kcvn := &KCPVconn{hc: session}
	err = kcvn.ApplyConf()
	if err != nil {
		return nil, err
	}
	return kcvn, nil
}

func init() {
	internet.KCPDialer = DialKCP
}
