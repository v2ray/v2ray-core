package kcp

import "net"

type UnreusableConnection struct {
	net.Conn
}

func (c UnreusableConnection) Reusable() bool {
	return false
}

func (c UnreusableConnection) SetReusable(bool) {}
