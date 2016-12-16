package scenarios

import (
	"github.com/golang/protobuf/proto"
	"sync/atomic"
	"time"
	"v2ray.com/core"
	v2net "v2ray.com/core/common/net"
)

var (
	port uint32 = 50000
)

func pickPort() v2net.Port {
	return v2net.Port(atomic.AddUint32(&port, 1))
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
