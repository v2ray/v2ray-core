package conf

import (
	"v2ray.com/core/common/serial"
)

type Buildable interface {
	Build() (*serial.TypedMessage, error)
}
