package web

import (
	"v2ray.com/core/app"
	"v2ray.com/core/common"
)

const (
	APP_ID = app.ID(8)
)

type Server interface {
	common.Releasable
	Handle()
}
