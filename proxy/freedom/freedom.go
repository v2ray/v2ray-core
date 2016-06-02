package freedom

import (
	"io"
	"net"
	"sync"

	"github.com/v2ray/v2ray-core/app"
	"github.com/v2ray/v2ray-core/app/dns"
	"github.com/v2ray/v2ray-core/common/alloc"
	"github.com/v2ray/v2ray-core/common/dice"
	v2io "github.com/v2ray/v2ray-core/common/io"
	"github.com/v2ray/v2ray-core/common/log"
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/common/retry"
	"github.com/v2ray/v2ray-core/proxy"
	"github.com/v2ray/v2ray-core/proxy/internal"
	"github.com/v2ray/v2ray-core/transport/hub"
	"github.com/v2ray/v2ray-core/transport/ray"
)

type FreedomConnection struct {
	domainStrategy DomainStrategy
	timeout        uint32
	dns            dns.Server
}

func NewFreedomConnection(config *Config, space app.Space) *FreedomConnection {
	f := &FreedomConnection{
		domainStrategy: config.DomainStrategy,
		timeout:        config.Timeout,
	}
	space.InitializeApplication(func() error {
		if config.DomainStrategy == DomainStrategyUseIP {
			if !space.HasApp(dns.APP_ID) {
				log.Error("Freedom: DNS server is not found in the space.")
				return app.ErrorMissingApplication
			}
			f.dns = space.GetApp(dns.APP_ID).(dns.Server)
		}
		return nil
	})
	return f
}

// @Private
func (this *FreedomConnection) ResolveIP(destination v2net.Destination) v2net.Destination {
	if !destination.Address().IsDomain() {
		return destination
	}

	ips := this.dns.Get(destination.Address().Domain())
	if len(ips) == 0 {
		log.Info("Freedom: DNS returns nil answer. Keep domain as is.")
		return destination
	}

	ip := ips[dice.Roll(len(ips))]
	var newDest v2net.Destination
	if destination.IsTCP() {
		newDest = v2net.TCPDestination(v2net.IPAddress(ip), destination.Port())
	} else {
		newDest = v2net.UDPDestination(v2net.IPAddress(ip), destination.Port())
	}
	log.Info("Freedom: Changing destination from ", destination, " to ", newDest)
	return newDest
}

func (this *FreedomConnection) Dispatch(destination v2net.Destination, payload *alloc.Buffer, ray ray.OutboundRay) error {
	log.Info("Freedom: Opening connection to ", destination)

	defer payload.Release()
	defer ray.OutboundInput().Release()
	defer ray.OutboundOutput().Close()

	var conn net.Conn
	if this.domainStrategy == DomainStrategyUseIP && destination.Address().IsDomain() {
		destination = this.ResolveIP(destination)
	}
	err := retry.Timed(5, 100).On(func() error {
		rawConn, err := hub.DialWithoutCache(destination)
		if err != nil {
			return err
		}
		conn = rawConn
		return nil
	})
	if err != nil {
		log.Error("Freedom: Failed to open connection to ", destination, ": ", err)
		return err
	}
	defer conn.Close()

	input := ray.OutboundInput()
	output := ray.OutboundOutput()
	var readMutex, writeMutex sync.Mutex
	readMutex.Lock()
	writeMutex.Lock()

	conn.Write(payload.Value)

	go func() {
		v2writer := v2io.NewAdaptiveWriter(conn)
		defer v2writer.Release()

		v2io.Pipe(input, v2writer)
		writeMutex.Unlock()
	}()

	go func() {
		defer readMutex.Unlock()

		var reader io.Reader = conn

		timeout := this.timeout
		if destination.IsUDP() {
			timeout = 16
		}
		if timeout > 0 {
			reader = v2net.NewTimeOutReader(int(timeout) /* seconds */, conn)
		}

		v2reader := v2io.NewAdaptiveReader(reader)
		defer v2reader.Release()

		v2io.Pipe(v2reader, output)
		ray.OutboundOutput().Close()
	}()

	writeMutex.Lock()
	if tcpConn, ok := conn.(*net.TCPConn); ok {
		tcpConn.CloseWrite()
	}
	readMutex.Lock()

	return nil
}

func init() {
	internal.MustRegisterOutboundHandlerCreator("freedom",
		func(space app.Space, config interface{}) (proxy.OutboundHandler, error) {
			return NewFreedomConnection(config.(*Config), space), nil
		})
}
