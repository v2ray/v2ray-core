package mux

import "v2ray.com/core/common/net"

const (
	maxParallel = 8
	maxTotal    = 128
)

type mergerWorker struct {
}

func (w *mergerWorker) isFull() bool {
	return true
}

type Merger struct {
	sessions map[net.Destination]mergerWorker
}
