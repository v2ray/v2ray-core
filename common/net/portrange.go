package net

type PortRange interface {
	From() uint16
	To() uint16
}
