package kcp

import (
	"errors"
	"net"
	"time"

	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/transport/internet"

	"github.com/xtaci/kcp-go"
)

type KCPVlistener struct {
	lst                    *kcp.Listener
	previousSocketid       map[int]uint32
	previousSocketid_mapid int
}

/*Accept Accept a KCP connection
Since KCP is stateless, if package deliver after it was closed,
It could be reconized as a new connection and call accept.
If we can detect that the connection is of such a kind,
we will discard that conn.
*/
func (kvl *KCPVlistener) Accept() (internet.Connection, error) {
	conn, err := kvl.lst.Accept()
	if err != nil {
		return nil, err
	}

	if kvl.previousSocketid == nil {
		kvl.previousSocketid = make(map[int]uint32)
	}

	var badbit bool = false

	for _, key := range kvl.previousSocketid {
		if key == conn.GetConv() {
			badbit = true
		}
	}
	if badbit {
		conn.Close()
		return nil, errors.New("KCP:ConnDup, Don't worry~")
	} else {
		kvl.previousSocketid_mapid++
		kvl.previousSocketid[kvl.previousSocketid_mapid] = conn.GetConv()
		/*
			Here we assume that count(connection) < 512
			This won't always true.
			More work might be necessary to deal with this in a better way.
		*/
		if kvl.previousSocketid_mapid >= 512 {
			delete(kvl.previousSocketid, kvl.previousSocketid_mapid-512)
		}
	}

	kcv := &KCPVconn{hc: conn}
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
	conntokeep time.Time
}

//var counter int

func (kcpvc *KCPVconn) Read(b []byte) (int, error) {
	ifb := time.Now().Add(time.Duration(effectiveConfig.ReadTimeout) * time.Second)
	if ifb.After(kcpvc.conntokeep) {
		kcpvc.conntokeep = ifb
	}
	kcpvc.hc.SetDeadline(kcpvc.conntokeep)
	return kcpvc.hc.Read(b)
}

func (kcpvc *KCPVconn) Write(b []byte) (int, error) {
	ifb := time.Now().Add(time.Duration(effectiveConfig.WriteTimeout) * time.Second)
	if ifb.After(kcpvc.conntokeep) {
		kcpvc.conntokeep = ifb
	}
	kcpvc.hc.SetDeadline(kcpvc.conntokeep)
	return kcpvc.hc.Write(b)
}

/*ApplyConf will apply kcpvc.conf to current Socket

It is recommmanded to call this func once and only once
*/
func (kcpvc *KCPVconn) ApplyConf() error {
	nodelay, interval, resend, nc := 0, 40, 0, 0
	switch effectiveConfig.Mode {
	case "normal":
		nodelay, interval, resend, nc = 0, 30, 2, 1
	case "fast":
		nodelay, interval, resend, nc = 0, 20, 2, 1
	case "fast2":
		nodelay, interval, resend, nc = 1, 20, 2, 1
	case "fast3":
		nodelay, interval, resend, nc = 1, 10, 2, 1
	}

	kcpvc.hc.SetNoDelay(nodelay, interval, resend, nc)
	kcpvc.hc.SetWindowSize(effectiveConfig.Sndwnd, effectiveConfig.Rcvwnd)
	kcpvc.hc.SetMtu(effectiveConfig.Mtu)
	kcpvc.hc.SetACKNoDelay(effectiveConfig.Acknodelay)
	kcpvc.hc.SetDSCP(effectiveConfig.Dscp)
	//counter++
	//log.Info(counter)
	return nil
}

/*Close Close the current conn
We have to delay the close of Socket for a few second
or the VMess EOF can be too late to send.
*/
func (kcpvc *KCPVconn) Close() error {
	go func() {
		time.Sleep(2000 * time.Millisecond)
		//counter--
		//log.Info(counter)
		kcpvc.hc.Close()
	}()
	return nil
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

func (this *KCPVconn) Reusable() bool {
	return false
}

func (this *KCPVconn) SetReusable(b bool) {

}

func ListenKCP(address v2net.Address, port v2net.Port) (internet.Listener, error) {
	laddr := address.String() + ":" + port.String()
	crypt, _ := kcp.NewNoneBlockCrypt(nil)
	kcl, err := kcp.ListenWithOptions(effectiveConfig.Fec, laddr, crypt)
	kcvl := &KCPVlistener{lst: kcl}
	return kcvl, err
}

func init() {
	internet.KCPListenFunc = ListenKCP
}
