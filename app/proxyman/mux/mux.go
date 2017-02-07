package mux

import "v2ray.com/core/common/net"

type mergerWorker struct {
}

type Merger struct {
	sessions map[net.Destination]mergerWorker
}
