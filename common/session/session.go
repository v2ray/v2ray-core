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

func ExportIDToError(ctx context.Context) errors.ExportOption {
	id := IDFromContext(ctx)
	return func(h *errors.ExportOptionHolder) {
		h.SessionID = uint32(id)
	}
}

type Inbound struct {
	Source  net.Destination
	Gateway net.Destination
	Tag     string
	// User is the user that authencates for the inbound. May be nil if the protocol allows anounymous traffic.
	User *protocol.MemoryUser
}

type Outbound struct {
	Target      net.Destination
	Gateway     net.Address
	ResolvedIPs []net.IP
}
