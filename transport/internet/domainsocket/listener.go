package domainsocket

import (
	"context"
	"net"
)

type Listener struct {
	ln net.Listener
}

func ListenDS(ctx context.Context, path string) (*Listener, error) {
	addr := new(net.UnixAddr)
	addr.Name = path
	addr.Net = "unix"
	li, err := net.ListenUnix("unix", addr)
	if err != nil {
		return nil, err
	}
	vln := &Listener{ln: li}
	return vln, nil
}

func (ls *Listener) Down() {
	ls.ln.Close()
}
