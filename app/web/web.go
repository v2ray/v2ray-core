package web

import "v2ray.com/core/app"

const (
	APP_ID = app.ID(8)
)

type WebServer interface {
	Handle()
}
