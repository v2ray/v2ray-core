package scenarios

import (
	"sync/atomic"
	"time"

	"net"

	"github.com/golang/protobuf/proto"
	"v2ray.com/core"
	v2net "v2ray.com/core/common/net"
)

var (
	port uint32 = 50000
)

func pickPort() v2net.Port {
	return v2net.Port(atomic.AddUint32(&port, 1))
}

func xor(b []byte) []byte {
	r := make([]byte, len(b))
	for i, v := range b {
		r[i] = v ^ 'c'
	}
	return r
}

func readFrom(conn net.Conn, timeout time.Duration, length int) []byte {
	b := make([]byte, 2048)
	totalBytes := 0
	deadline := time.Now().Add(timeout)
	conn.SetReadDeadline(deadline)
	for totalBytes < length {
		if time.Now().After(deadline) {
			break
		}
		n, err := conn.Read(b[totalBytes:])
		if err != nil {
			break
		}
		totalBytes += n
	}
	return b[:totalBytes]
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
