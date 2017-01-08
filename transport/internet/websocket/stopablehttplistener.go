package websocket

import (
	"net"

	"v2ray.com/core/common/errors"
)

type StoppableListener struct {
	net.Listener //Wrapped listener
}

func NewStoppableListener(l net.Listener) (*StoppableListener, error) {

	retval := &StoppableListener{}
	retval.Listener = l
	return retval, nil
}

var StoppedError = errors.New("Listener stopped")

func (sl *StoppableListener) Accept() (net.Conn, error) {
	newConn, err := sl.Listener.Accept()
	return newConn, err

}

func (sl *StoppableListener) Stop() {
	sl.Listener.Close()
}
