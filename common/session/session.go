// Package session provides functions for sessions of incoming requests.
package session // import "v2ray.com/core/common/session"

import (
	"context"
	"math/rand"

	"v2ray.com/core/common/errors"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/protocol"
)

// ID of a session.
type ID uint32

// NewID generates a new ID. The generated ID is high likely to be unique, but not cryptographically secure.
// The generated ID will never be 0.
func NewID() ID {
	for {
		id := ID(rand.Uint32())
		if id != 0 {
			return id
		}
	}
}

// ExportIDToError transfers session.ID into an error object, for logging purpose.
// This can be used with error.WriteToLog().
func ExportIDToError(ctx context.Context) errors.ExportOption {
	id := IDFromContext(ctx)
	return func(h *errors.ExportOptionHolder) {
		h.SessionID = uint32(id)
	}
}

// Inbound is the metadata of an inbound connection.
type Inbound struct {
	// Source address of the inbound connection.
	Source net.Destination
	// Getaway address
	Gateway net.Destination
	// Tag of the inbound proxy that handles the connection.
	Tag string
	// User is the user that authencates for the inbound. May be nil if the protocol allows anounymous traffic.
	User *protocol.MemoryUser
}

// Outbound is the metadata of an outbound connection.
type Outbound struct {
	// Target address of the outbound connection.
	Target net.Destination
	// Gateway address
	Gateway net.Address
	// ResolvedIPs is the resolved IP addresses, if the Targe is a domain address.
	ResolvedIPs []net.IP
}

type SniffingRequest struct {
	OverrideDestinationForProtocol []string
	Enabled                        bool
}

// Content is the metadata of the connection content.
type Content struct {
	// Protocol of current content.
	Protocol string

	SniffingRequest SniffingRequest

	Attributes map[string]interface{}

	SkipRoutePick bool
}

func (c *Content) SetAttribute(name string, value interface{}) {
	if c.Attributes == nil {
		c.Attributes = make(map[string]interface{})
	}
	c.Attributes[name] = value
}

func (c *Content) Attribute(name string) interface{} {
	if c.Attributes == nil {
		return nil
	}
	return c.Attributes[name]
}
