package hub

import (
	"errors"
	"net"
	"time"

	"github.com/v2ray/v2ray-core/common/log"
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/transport"
	"github.com/v2ray/v2ray-core/transport/hub/kcpv"
	"github.com/xtaci/kcp-go"
)

type KCPVlistener struct {
	lst                    *kcp.Listener
	conf                   *kcpv.Config
	previousSocketid       map[int]uint32
	previousSocketid_mapid int
}

/*Accept Accept a KCP connection
Since KCP is stateless, if package deliver after it was closed,
It could be reconized as a new connection and call accept.
If we can detect that the connection is of such a kind,
we will discard that conn.
*/
func (kvl *KCPVlistener) Accept() (net.Conn, error) {
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
		if kvl.previousSocketid_mapid >= 512 {
			delete(kvl.previousSocketid, kvl.previousSocketid_mapid-512)
		}
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

//var counter int

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
	//counter++
	//log.Info(counter)
	return nil
}

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

func DialKCP(dest v2net.Destination) (*KCPVconn, error) {
	kcpconf := transport.KcpConfig
	cpip, _ := kcpv.GetChipher(kcpconf.Key)
	kcv, err := kcp.DialWithOptions(kcpconf.AdvancedConfigs.Fec, dest.NetAddr(), cpip)
	if err != nil {
		return nil, err
	}
	kcvn := &KCPVconn{hc: kcv}
	kcvn.conf = kcpconf
	err = kcvn.ApplyConf()
	if err != nil {
		return nil, err
	}
	return kcvn, nil
}

func ListenKCP(address v2net.Address, port v2net.Port) (*KCPVlistener, error) {
	kcpconf := transport.KcpConfig
	cpip, _ := kcpv.GetChipher(kcpconf.Key)
	laddr := address.String() + ":" + port.String()
	kcl, err := kcp.ListenWithOptions(kcpconf.AdvancedConfigs.Fec, laddr, cpip)
	kcvl := &KCPVlistener{lst: kcl, conf: kcpconf}
	return kcvl, err
}
