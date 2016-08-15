package ws

import "crypto/tls"

func getstopableTLSlistener(cert, key, listenaddr string) (*StoppableListener, error) {
	cer, err := tls.LoadX509KeyPair(cert, key)
	if err != nil {
		return nil, err
	}

	config := &tls.Config{Certificates: []tls.Certificate{cer}}
	ln, err := tls.Listen("tcp", listenaddr, config)
	if err != nil {
		return nil, err
	}
	lns, err := NewStoppableListener(ln)
	return lns, err
}
