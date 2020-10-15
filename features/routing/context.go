package routing

import (
	"v2ray.com/core/common/net"
)

// Context is a feature to store connection information for routing.
//
// v2ray:api:stable
type Context interface {
	// GetInboundTag returns the tag of the inbound the connection was from.
	GetInboundTag() string

	// GetSourcesIPs returns the source IPs bound to the connection.
	GetSourceIPs() []net.IP

	// GetSourcePort returns the source port of the connection.
	GetSourcePort() net.Port

	// GetTargetIPs returns the target IP of the connection or resolved IPs of target domain.
	GetTargetIPs() []net.IP

	// GetTargetPort returns the target port of the connection.
	GetTargetPort() net.Port

	// GetTargetDomain returns the target domain of the connection, if exists.
	GetTargetDomain() string

	// GetNetwork returns the network type of the connection.
	GetNetwork() net.Network

	// GetProtocol returns the protocol from the connection content, if sniffed out.
	GetProtocol() string

	// GetUser returns the user email from the connection content, if exists.
	GetUser() string

	// GetAttributes returns extra attributes from the conneciont content.
	GetAttributes() map[string]string
}
