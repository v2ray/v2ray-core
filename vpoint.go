package core

import (
	"fmt"
)

type VPoint struct {
	config      VConfig
	connHandler ConnectionHandler
}

func NewVPoint(config *VConfig) (*VPoint, error) {
	var vpoint *VPoint
	return vpoint, nil
}

type ConnectionHandler interface {
	Listen(port uint16) error
}

func (vp *VPoint) Start() error {
	if vp.config.Port <= 0 {
		return fmt.Errorf("Invalid port %d", vp.config.Port)
	}
	vp.connHandler.Listen(vp.config.Port)
	return nil
}
