package tcp

import (
	"context"
	"fmt"
	"io"

	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/task"
	"v2ray.com/core/transport/internet"
	"v2ray.com/core/transport/pipe"
)

type Server struct {
	Port         net.Port
	MsgProcessor func(msg []byte) []byte
	ShouldClose  bool
	SendFirst    []byte
	Listen       net.Address
	listener     net.Listener
}

func (server *Server) Start() (net.Destination, error) {
	return server.StartContext(context.Background())
}

func (server *Server) StartContext(ctx context.Context) (net.Destination, error) {
	listenerAddr := server.Listen
	if listenerAddr == nil {
		listenerAddr = net.LocalHostIP
	}
	listener, err := internet.ListenSystem(ctx, &net.TCPAddr{
		IP:   listenerAddr.IP(),
		Port: int(server.Port),
	})
	if err != nil {
		return net.Destination{}, err
	}

	localAddr := listener.Addr().(*net.TCPAddr)
	server.Port = net.Port(localAddr.Port)
	server.listener = listener
	go server.acceptConnections(listener.(*net.TCPListener))

	return net.TCPDestination(net.IPAddress(localAddr.IP), net.Port(localAddr.Port)), nil
}

func (server *Server) acceptConnections(listener *net.TCPListener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Failed accept TCP connection: %v\n", err)
			return
		}

		go server.handleConnection(conn)
	}
}

func (server *Server) handleConnection(conn net.Conn) {
	if len(server.SendFirst) > 0 {
		conn.Write(server.SendFirst)
	}

	pReader, pWriter := pipe.New(pipe.WithoutSizeLimit())
	err := task.Run(task.Parallel(func() error {
		defer pWriter.Close() // nolint: errcheck

		for {
			b := buf.New()
			if _, err := b.ReadFrom(conn); err != nil {
				if err == io.EOF {
					return nil
				}
				return err
			}
			copy(b.Bytes(), server.MsgProcessor(b.Bytes()))
			if err := pWriter.WriteMultiBuffer(buf.MultiBuffer{b}); err != nil {
				return err
			}
		}
	}, func() error {
		defer pReader.CloseError()

		w := buf.NewWriter(conn)
		for {
			mb, err := pReader.ReadMultiBuffer()
			if err != nil {
				if err == io.EOF {
					return nil
				}
				return err
			}
			if err := w.WriteMultiBuffer(mb); err != nil {
				return err
			}
		}
	}))()

	if err != nil {
		fmt.Println("failed to transfer data: ", err.Error())
	}

	conn.Close() // nolint: errcheck
}

func (server *Server) Close() error {
	return server.listener.Close()
}
