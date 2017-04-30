package unix

import (
	"context"
	"net"
	"os"
	"sync"

	"golang.org/x/sys/unix"
	"v2ray.com/core/app/vpndialer"
	"v2ray.com/core/common"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/common/serial"
	"v2ray.com/core/transport/internet"
)

//go:generate go run $GOPATH/src/v2ray.com/core/tools/generrorgen/main.go -pkg unix -path App,VPNDialer,Unix

type status int

const (
	statusNew status = iota
	statusOK
	statusFail
)

type fdStatus struct {
	status   status
	fd       int
	callback chan<- error
}

type protector struct {
	sync.Mutex
	address string
	conn    *net.UnixConn
	status  chan fdStatus
}

func readFrom(conn *net.UnixConn, schan chan<- fdStatus) {
	var payload [6]byte
	for {
		_, err := conn.Read(payload[:])
		if err != nil {
			break
		}
		fd := serial.BytesToInt(payload[1:5])
		s := status(payload[5])
		schan <- fdStatus{
			fd:     fd,
			status: s,
		}
	}
}

func (m *protector) dial() (*net.UnixConn, error) {
	m.Lock()
	defer m.Unlock()

	if m.conn != nil {
		return m.conn, nil
	}

	conn, err := net.DialUnix("unix", nil, &net.UnixAddr{
		Name: m.address,
		Net:  "unix",
	})
	if err != nil {
		return nil, err
	}
	m.conn = conn
	m.status = make(chan fdStatus, 32)
	go readFrom(conn, m.status)
	go m.monitor(conn)
	return conn, nil
}

func (m *protector) close() {
	m.Lock()
	defer m.Unlock()
	if m.conn == nil {
		return
	}
	m.conn.Close()
	m.conn = nil
}

func (m *protector) monitor(c *net.UnixConn) {
	pendingFd := make(map[int]chan<- error, 32)
	for s := range m.status {
		switch s.status {
		case statusNew:
			pendingFd[s.fd] = s.callback
		case statusOK:
			if c, f := pendingFd[s.fd]; f {
				close(c)
				delete(pendingFd, s.fd)
			}
		case statusFail:
			if c, f := pendingFd[s.fd]; f {
				c <- newError("failed to protect fd")
				close(c)
				delete(pendingFd, s.fd)
			}
		}
	}
}

func (m *protector) protect(fd int) error {
	conn, err := m.dial()
	if err != nil {
		return err
	}

	var payload [6]byte
	serial.IntToBytes(fd, payload[1:1])
	payload[5] = byte(statusNew)
	if _, err := conn.Write(payload[:]); err != nil {
		return err
	}

	wait := make(chan error)
	m.status <- fdStatus{
		status:   statusNew,
		fd:       fd,
		callback: wait,
	}
	return <-wait
}

type App struct {
	protector *protector
	dialer    *Dialer
}

func NewApp(ctx context.Context, config *vpndialer.Config) (*App, error) {
	a := &App{
		dialer: &Dialer{},
		protector: &protector{
			address: config.Address,
		},
	}
	a.dialer.protect = a.protector.protect
	return a, nil
}

func (*App) Interface() interface{} {
	return (*App)(nil)
}

func (a *App) Start() error {
	internet.UseAlternativeSystemDialer(a.dialer)
	return nil
}

func (a *App) Close() {
	internet.UseAlternativeSystemDialer(nil)
}

type Dialer struct {
	protect func(fd int) error
}

func socket(dest v2net.Destination) (int, error) {
	switch dest.Network {
	case v2net.Network_TCP:
		return unix.Socket(unix.AF_INET6, unix.SOCK_STREAM, unix.IPPROTO_TCP)
	case v2net.Network_UDP:
		return unix.Socket(unix.AF_INET6, unix.SOCK_DGRAM, unix.IPPROTO_UDP)
	default:
		return 0, newError("unknown network ", dest.Network)
	}
}

func getIP(addr v2net.Address) (net.IP, error) {
	if addr.Family().Either(v2net.AddressFamilyIPv4, v2net.AddressFamilyIPv6) {
		return addr.IP(), nil
	}
	ips, err := net.LookupIP(addr.Domain())
	if err != nil {
		return nil, err
	}
	return ips[0], nil
}

func (d *Dialer) Dial(ctx context.Context, source v2net.Address, dest v2net.Destination) (net.Conn, error) {
	fd, err := socket(dest)
	if err != nil {
		return nil, err
	}
	if err := d.protect(fd); err != nil {
		return nil, err
	}

	ip, err := getIP(dest.Address)
	if err != nil {
		return nil, err
	}

	addr := &unix.SockaddrInet6{
		Port:   int(dest.Port),
		ZoneId: 0,
	}
	copy(addr.Addr[:], ip.To16())

	if err := unix.Connect(fd, addr); err != nil {
		return nil, err
	}

	file := os.NewFile(uintptr(fd), "Socket")
	return net.FileConn(file)
}

func init() {
	common.Must(common.RegisterConfig((*vpndialer.Config)(nil), func(ctx context.Context, config interface{}) (interface{}, error) {
		return NewApp(ctx, config.(*vpndialer.Config))
	}))
}
