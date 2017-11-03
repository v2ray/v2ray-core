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
	addr.Net = "unixpacket"
	li, err := net.ListenUnix("unixpacket", addr)
	if err != nil {
		return nil, err
	}
	vln := &Listener{ln: li}
	return vln, nil
}
