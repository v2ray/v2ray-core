package http

import (
	"net"
	"net/http"
	"strings"

	"github.com/v2ray/v2ray-core/app"
	_ "github.com/v2ray/v2ray-core/common/log"
	v2net "github.com/v2ray/v2ray-core/common/net"
)

type HttpProxyServer struct {
	accepting  bool
	dispatcher app.PacketDispatcher
	config     Config
}

func NewHttpProxyServer(dispatcher app.PacketDispatcher, config Config) *HttpProxyServer {
	return &HttpProxyServer{
		dispatcher: dispatcher,
		config:     config,
	}
}

func (this *HttpProxyServer) Listen(port v2net.Port) error {
	server := http.Server{
		Addr:    ":" + port.String(),
		Handler: this,
	}
	return server.ListenAndServe()
}

func (this *HttpProxyServer) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	if strings.ToUpper(request.Method) == "CONNECT" {
		host, port, err := net.SplitHostPort(request.URL.Host)
		if err != nil {
			if strings.Contains(err.(*net.AddrError).Err, "missing port") {
				host = request.URL.Host
				port = "80"
			} else {
				http.Error(w, "Bad Request", 400)
				return
			}
		}
		_ = host + port
	} else {

	}
}

func (this *HttpProxyServer) handleConnect(response http.ResponseWriter, request *http.Request) {

}
