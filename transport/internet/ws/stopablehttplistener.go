package ws

import (
	"errors"
	"net"
)

type StoppableListener struct {
	net.Listener          //Wrapped listener
	stop         chan int //Channel used only to indicate listener should shutdown
}

func NewStoppableListener(l net.Listener) (*StoppableListener, error) {
	/*
		tcpL, ok := l.(*net.TCPListener)

		if !ok {
			return nil, errors.New("Cannot wrap listener")
		}
	*/
	retval := &StoppableListener{}
	retval.Listener = l
	retval.stop = make(chan int)

	return retval, nil
}

var StoppedError = errors.New("Listener stopped")

func (sl *StoppableListener) Accept() (net.Conn, error) {

	for {
		newConn, err := sl.Listener.Accept()

		//Check for the channel being closed
		select {
		case <-sl.stop:
			return nil, StoppedError
		default:
			//If the channel is still open, continue as normal
		}

		if err != nil {
			netErr, ok := err.(net.Error)

			//If this is a timeout, then continue to wait for
			//new connections
			if ok && netErr.Timeout() && netErr.Temporary() {
				continue
			}
		}

		return newConn, err
	}
}

func (sl *StoppableListener) Stop() {
	close(sl.stop)
}
