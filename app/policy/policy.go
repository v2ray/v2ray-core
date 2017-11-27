package policy

import (
	"v2ray.com/core/app"
)

type Interface interface {
	GetPolicy(level uint32) Policy
}

func PolicyFromSpace(space app.Space) Interface {
	app := space.GetApplication((*Interface)(nil))
	if app == nil {
		return nil
	}
	return app.(Interface)
}
