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

	kcv := &KCPVconn{hc: conn}
	kcv.conf = kvl.conf
	err = kcv.ApplyConf()
	if err != nil {
		return nil, err
	}
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
func (kcpvc *KCPVconn) ApplyConf() error {
	nodelay, interval, resend, nc := 0, 40, 0, 0
	if kcpvc.conf.Mode != "manual" {
		switch kcpvc.conf.Mode {
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
		log.Error("kcp: Failed to Apply configure: Manual mode is not supported.(yet!)")
		return errors.New("kcp: Manual Not Implemented")
	}

	kcpvc.hc.SetNoDelay(nodelay, interval, resend, nc)
	kcpvc.hc.SetWindowSize(kcpvc.conf.AdvancedConfigs.Sndwnd, kcpvc.conf.AdvancedConfigs.Rcvwnd)
	kcpvc.hc.SetMtu(kcpvc.conf.AdvancedConfigs.Mtu)
	kcpvc.hc.SetACKNoDelay(kcpvc.conf.AdvancedConfigs.Acknodelay)
	kcpvc.hc.SetDSCP(kcpvc.conf.AdvancedConfigs.Dscp)
	return nil
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
