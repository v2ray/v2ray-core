package policy

import (
	"v2ray.com/core/app"
)

// Manager is an utility to manage policy per user level.
type Manager interface {
	// GetPolicy returns the Policy for the given user level.
	GetPolicy(level uint32) Policy
}

// FromSpace returns the policy.Manager in a space.
func FromSpace(space app.Space) Manager {
	app := space.GetApplication((*Manager)(nil))
	if app == nil {
		return nil
	}
	return app.(Manager)
}
