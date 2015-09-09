package core

import (
	"fmt"
)

// VPoint is an single server in V2Ray system.
type VPoint struct {
	config     VConfig
	ichFactory InboundConnectionHandlerFactory
	ochFactory OutboundConnectionHandlerFactory
}

// NewVPoint returns a new VPoint server based on given configuration.
// The server is not started at this point.
func NewVPoint(config *VConfig) (*VPoint, error) {
	var vpoint = new(VPoint)
	vpoint.config = *config
	return vpoint, nil
}

type InboundConnectionHandlerFactory interface {
	Create(vPoint *VPoint) (InboundConnectionHandler, error)
}

type InboundConnectionHandler interface {
	Listen(port uint16) error
}

type OutboundConnectionHandlerFactory interface {
	Create(vPoint *VPoint) (OutboundConnectionHandler, error)
}

type OutboundConnectionHandler interface {
	Start(vray *OutboundVRay) error
}

// Start starts the VPoint server, and return any error during the process.
// In the case of any errors, the state of the server is unpredicatable.
func (vp *VPoint) Start() error {
	if vp.config.Port <= 0 {
		return fmt.Errorf("Invalid port %d", vp.config.Port)
	}
	inboundConnectionHandler, err := vp.ichFactory.Create(vp)
	if err != nil {
		return err
	}
	err = inboundConnectionHandler.Listen(vp.config.Port)
	return nil
}
