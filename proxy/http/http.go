package http

import (
	"net"
	// "net/http"

	"github.com/v2ray/v2ray-core/app"
	"github.com/v2ray/v2ray-core/common/log"
	v2net "github.com/v2ray/v2ray-core/common/net"
	jsonconfig "github.com/v2ray/v2ray-core/proxy/http/config/json"
)

type HttpProxyServer struct {
	accepting  bool
	dispatcher app.PacketDispatcher
	config     *jsonconfig.HttpProxyConfig
}

func NewHttpProxyServer(dispatcher app.PacketDispatcher, config *jsonconfig.HttpProxyConfig) *HttpProxyServer {
	return &HttpProxyServer{
		dispatcher: dispatcher,
		config:     config,
	}
}

func (server *HttpProxyServer) Listen(port v2net.Port) error {
	_, err := net.ListenTCP("tcp", &net.TCPAddr{
		IP:   []byte{0, 0, 0, 0},
		Port: int(port),
		Zone: "",
	})
	if err != nil {
		log.Error("HTTP Proxy failed to listen on port %d: %v", port, err)
		return err
	}
	return nil
}
