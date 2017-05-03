package scenarios

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
	"v2ray.com/core"
	"v2ray.com/core/app/log"
	"v2ray.com/core/common"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/common/retry"
)

func pickPort() v2net.Port {
	listener, err := net.Listen("tcp4", ":0")
	common.Must(err)
	defer listener.Close()

	addr := listener.Addr().(*net.TCPAddr)
	return v2net.Port(addr.Port)
}

func pickUDPPort() v2net.Port {
	conn, err := net.ListenUDP("udp4", &net.UDPAddr{
		IP:   v2net.LocalHostIP.IP(),
		Port: 0,
	})
	common.Must(err)
	defer conn.Close()

	addr := conn.LocalAddr().(*net.UDPAddr)
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
	n, err := io.ReadFull(conn, b[:length])
	if err != nil {
		fmt.Println("Unexpected error from readFrom:", err)
	}
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

	if err := proc.Start(); err != nil {
		return err
	}

	time.Sleep(time.Second)

	runningServers = append(runningServers, proc)

	return nil
}

var (
	runningServers    = make([]*exec.Cmd, 0, 10)
	testBinaryPath    string
	testBinaryPathGen sync.Once
)

func genTestBinaryPath() {
	testBinaryPathGen.Do(func() {
		var tempDir string
		common.Must(retry.Timed(5, 100).On(func() error {
			dir, err := ioutil.TempDir("", "v2ray")
			if err != nil {
				return err
			}
			tempDir = dir
			return nil
		}))
		file := filepath.Join(tempDir, "v2ray.test")
		if runtime.GOOS == "windows" {
			file += ".exe"
		}
		testBinaryPath = file
		fmt.Printf("Generated binary path: %s\n", file)
	})
}

func GetSourcePath() string {
	return filepath.Join("v2ray.com", "core", "main")
}

func CloseAllServers() {
	log.Trace(errors.New("Closing all servers."))
	for _, server := range runningServers {
		server.Process.Signal(os.Interrupt)
	}
	for _, server := range runningServers {
		server.Process.Wait()
	}
	runningServers = make([]*exec.Cmd, 0, 10)
	log.Trace(errors.New("All server closed."))
}
