package scenarios

import (
	"io"
	"net"
	"time"

	"github.com/golang/protobuf/proto"
	"v2ray.com/core"
	"v2ray.com/core/common"
	v2net "v2ray.com/core/common/net"
)

func pickPort() v2net.Port {
	listener, err := net.Listen("tcp4", ":0")
	common.Must(err)
	defer listener.Close()

	addr := listener.Addr().(*net.TCPAddr)
	return v2net.Port(addr.Port)
}

func xor(b []byte) []byte {
	r := make([]byte, len(b))
	for i, v := range b {
		r[i] = v ^ 'c'
	}
	return r
}

func readFrom(conn net.Conn, timeout time.Duration, length int) []byte {
	b := make([]byte, length)
	deadline := time.Now().Add(timeout)
	conn.SetReadDeadline(deadline)
	n, _ := io.ReadFull(conn, b[:length])
	return b[:n]
}

func InitializeServerConfig(config *core.Config) error {
	err := BuildV2Ray()
	if err != nil {
		return err
	}

	configBytes, err := proto.Marshal(config)
	if err != nil {
		return err
	}
	proc := RunV2RayProtobuf(configBytes)

	err = proc.Start()
	if err != nil {
		return err
	}

	time.Sleep(time.Second)

	runningServers = append(runningServers, proc)

	return nil
}
