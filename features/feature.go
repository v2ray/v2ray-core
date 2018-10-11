package features

import "v2ray.com/core/common"

// Feature is the interface for V2Ray features. All features must implement this interface.
// All existing features have an implementation in app directory. These features can be replaced by third-party ones.
type Feature interface {
	common.Runnable
}
