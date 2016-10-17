package conf

import (
	"v2ray.com/core/common/loader"
)

type Buildable interface {
	Build() (*loader.TypedSettings, error)
}
