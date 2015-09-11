package core

import (
	"fmt"

	v2net "github.com/v2ray/v2ray-core/net"
)

// VPoint is an single server in V2Ray system.
type VPoint struct {
	Config     VConfig
	ichFactory InboundConnectionHandlerFactory
	ochFactory OutboundConnectionHandlerFactory
}

// NewVPoint returns a new VPoint server based on given configuration.
// The server is not started at this point.
func NewVPoint(config *VConfig, ichFactory InboundConnectionHandlerFactory, ochFactory OutboundConnectionHandlerFactory) (*VPoint, error) {
	var vpoint = new(VPoint)
	vpoint.Config = *config
	vpoint.ichFactory = ichFactory
	vpoint.ochFactory = ochFactory

	return vpoint, nil
}

type InboundConnectionHandlerFactory interface {
	Create(vPoint *VPoint) (InboundConnectionHandler, error)
}

type InboundConnectionHandler interface {
	Listen(port uint16) error
}

type OutboundConnectionHandlerFactory interface {
	Create(vPoint *VPoint, dest v2net.VAddress) (OutboundConnectionHandler, error)
}

type OutboundConnectionHandler interface {
	Start(vray OutboundVRay) error
}

// Start starts the VPoint server, and return any error during the process.
// In the case of any errors, the state of the server is unpredicatable.
func (vp *VPoint) Start() error {
	if vp.Config.Port <= 0 {
		return fmt.Errorf("Invalid port %d", vp.Config.Port)
	}
	inboundConnectionHandler, err := vp.ichFactory.Create(vp)
	if err != nil {
		return err
	}
	err = inboundConnectionHandler.Listen(vp.Config.Port)
	return nil
}

func (vp *VPoint) NewInboundConnectionAccepted(destination v2net.VAddress) InboundVRay {
	ray := NewVRay()
	// TODO: handle error
	och, _ := vp.ochFactory.Create(vp, destination)
	_ = och.Start(ray)
	return ray
}
