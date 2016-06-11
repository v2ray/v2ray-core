package hub

import (
	"errors"
	"net"
	"time"

	"github.com/v2ray/v2ray-core/common/log"
	"github.com/v2ray/v2ray-core/transport/hub/kcpv"
	"github.com/xtaci/kcp-go"
)

type KCPVlistener struct {
	lst  *kcp.Listener
	conf *kcpv.Config
}

func (kvl *KCPVlistener) Accept() (*KCPVconn, error) {
	conn, err := kvl.lst.Accept()
	if err != nil {
		return nil, err
	}
	nodelay, interval, resend, nc := 0, 40, 0, 0
	if kvl.conf.Mode != "manual" {
		switch kvl.conf.Mode {
		case "normal":
			nodelay, interval, resend, nc = 0, 30, 2, 1
		case "fast":
			nodelay, interval, resend, nc = 0, 20, 2, 1
		case "fast2":
			nodelay, interval, resend, nc = 1, 20, 2, 1
		case "fast3":
			nodelay, interval, resend, nc = 1, 10, 2, 1
		}
	} else {
		log.Error("kcp: Accepted Unsuccessfully: Manual mode is not supported.(yet!)")
		return nil, errors.New("kcp: Manual Not Implemented")
	}

	conn.SetNoDelay(nodelay, interval, resend, nc)
	conn.SetWindowSize(kvl.conf.AdvancedConfigs.Sndwnd, kvl.conf.AdvancedConfigs.Rcvwnd)
	conn.SetMtu(kvl.conf.AdvancedConfigs.Mtu)
	conn.SetACKNoDelay(kvl.conf.AdvancedConfigs.Acknodelay)
	conn.SetDSCP(kvl.conf.AdvancedConfigs.Dscp)

	kcv := &KCPVconn{hc: conn}
	kcv.conf = kvl.conf
	return kcv, nil
}
func (kvl *KCPVlistener) Close() error {
	return kvl.lst.Close()
}
func (kvl *KCPVlistener) Addr() net.Addr {
	return kvl.lst.Addr()
}

type KCPVconn struct {
	hc         *kcp.UDPSession
	conf       *kcpv.Config
	conntokeep time.Time
}

func (kcpvc *KCPVconn) Read(b []byte) (int, error) {
	ifb := time.Now().Add(time.Duration(kcpvc.conf.AdvancedConfigs.ReadTimeout) * time.Second)
	if ifb.After(kcpvc.conntokeep) {
		kcpvc.conntokeep = ifb
	}
	kcpvc.hc.SetDeadline(kcpvc.conntokeep)
	return kcpvc.hc.Read(b)
}

func (kcpvc *KCPVconn) Write(b []byte) (int, error) {
	ifb := time.Now().Add(time.Duration(kcpvc.conf.AdvancedConfigs.WriteTimeout) * time.Second)
	if ifb.After(kcpvc.conntokeep) {
		kcpvc.conntokeep = ifb
	}
	kcpvc.hc.SetDeadline(kcpvc.conntokeep)
	return kcpvc.hc.Write(b)
}

func (kcpvc *KCPVconn) Close() error {

	return kcpvc.hc.Close()
}

func (kcpvc *KCPVconn) LocalAddr() net.Addr {
	return kcpvc.hc.LocalAddr()
}

func (kcpvc *KCPVconn) RemoteAddr() net.Addr {
	return kcpvc.hc.RemoteAddr()
}

func (kcpvc *KCPVconn) SetDeadline(t time.Time) error {
	return kcpvc.hc.SetDeadline(t)
}

func (kcpvc *KCPVconn) SetReadDeadline(t time.Time) error {
	return kcpvc.hc.SetReadDeadline(t)
}

func (kcpvc *KCPVconn) SetWriteDeadline(t time.Time) error {
	return kcpvc.hc.SetWriteDeadline(t)
}
