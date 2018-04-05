package domainsocket

import (
	"context"
	"net"
)

func DialDS(ctx context.Context, path string) (*net.UnixConn, error) {
	resolvedAddress, err := net.ResolveUnixAddr("unix", path)
	if err != nil {
		return nil, err
	}
	dialedUnix, err := net.DialUnix("unix", nil, resolvedAddress)
	if err != nil {
		return nil, err
	}
	return dialedUnix, nil
}
