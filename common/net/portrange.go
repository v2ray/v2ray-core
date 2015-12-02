package net

type PortRange interface {
	From() Port
	To() Port
}
