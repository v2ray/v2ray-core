package websocket

import "crypto/tls"

func getstopableTLSlistener(tlsConfig *tls.Config, listenaddr string) (*StoppableListener, error) {
	ln, err := tls.Listen("tcp", listenaddr, tlsConfig)
	if err != nil {
		return nil, err
	}
	lns, err := NewStoppableListener(ln)
	return lns, err
}
