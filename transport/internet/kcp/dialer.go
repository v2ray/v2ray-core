package kcp

import (
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/transport/internet"

	"github.com/xtaci/kcp-go"
)

func DialKCP(src v2net.Address, dest v2net.Destination) (internet.Connection, error) {
	cpip, _ := kcp.NewNoneBlockCrypt(nil)
	kcv, err := kcp.DialWithOptions(effectiveConfig.Fec, dest.NetAddr(), cpip)
	if err != nil {
		return nil, err
	}
	kcvn := &KCPVconn{hc: kcv}
	err = kcvn.ApplyConf()
	if err != nil {
		return nil, err
	}
	return kcvn, nil
}

func init() {
	internet.KCPDialer = DialKCP
}
