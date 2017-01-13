package core

import (
	"context"

	"v2ray.com/core/app"
	"v2ray.com/core/app/dispatcher"
	"v2ray.com/core/app/dns"
	proxydialer "v2ray.com/core/app/proxy"
	"v2ray.com/core/app/proxyman"
	"v2ray.com/core/common"
	"v2ray.com/core/common/log"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/proxy"
)

// Point shell of V2Ray.
type Point struct {
	inboundHandlers       []InboundDetourHandler
	taggedInboundHandlers map[string]InboundDetourHandler

	outboundHandlers       []proxy.OutboundHandler
	taggedOutboundHandlers map[string]proxy.OutboundHandler

	space app.Space
}

// NewPoint returns a new Point server based on given configuration.
// The server is not started at this point.
func NewPoint(pConfig *Config) (*Point, error) {
	var vpoint = new(Point)

	if err := pConfig.Transport.Apply(); err != nil {
		return nil, err
	}

	if err := pConfig.Log.Apply(); err != nil {
		return nil, err
	}

	space := app.NewSpace()
	ctx := app.ContextWithSpace(context.Background(), space)

	vpoint.space = space
	vpoint.space.AddApplication(vpoint)

	outboundHandlerManager := proxyman.OutboundHandlerManagerFromSpace(space)
	if outboundHandlerManager == nil {
		o, err := app.CreateAppFromConfig(ctx, new(proxyman.OutboundConfig))
		if err != nil {
			return nil, err
		}
		space.AddApplication(o)
		outboundHandlerManager = o.(proxyman.OutboundHandlerManager)
	}

	proxyDialer := proxydialer.OutboundProxyFromSpace(space)
	if proxyDialer == nil {
		p, err := app.CreateAppFromConfig(ctx, new(proxydialer.Config))
		if err != nil {
			return nil, err
		}
		space.AddApplication(p)
		proxyDialer = p.(*proxydialer.OutboundProxy)
	}
	proxyDialer.RegisterDialer()

	for _, appSettings := range pConfig.App {
		settings, err := appSettings.GetInstance()
		if err != nil {
			return nil, err
		}
		application, err := app.CreateAppFromConfig(ctx, settings)
		if err != nil {
			return nil, err
		}
		if err := space.AddApplication(application); err != nil {
			return nil, err
		}
	}

	dnsServer := dns.FromSpace(space)
	if dnsServer == nil {
		dnsConfig := &dns.Config{
			NameServers: []*v2net.Endpoint{{
				Address: v2net.NewIPOrDomain(v2net.LocalHostDomain),
			}},
		}
		d, err := app.CreateAppFromConfig(ctx, dnsConfig)
		if err != nil {
			return nil, err
		}
		space.AddApplication(d)
		dnsServer = d.(dns.Server)
	}

	disp := dispatcher.FromSpace(space)
	if disp == nil {
		d, err := app.CreateAppFromConfig(ctx, new(dispatcher.Config))
		if err != nil {
			return nil, err
		}
		space.AddApplication(d)
		disp = d.(dispatcher.Interface)
	}

	vpoint.inboundHandlers = make([]InboundDetourHandler, 0, 8)
	vpoint.taggedInboundHandlers = make(map[string]InboundDetourHandler)
	for _, inbound := range pConfig.Inbound {
		allocConfig := inbound.GetAllocationStrategyValue()
		var inboundHandler InboundDetourHandler
		switch allocConfig.Type {
		case AllocationStrategy_Always:
			dh, err := NewInboundDetourHandlerAlways(ctx, inbound)
			if err != nil {
				log.Error("V2Ray: Failed to create detour handler: ", err)
				return nil, common.ErrBadConfiguration
			}
			inboundHandler = dh
		case AllocationStrategy_Random:
			dh, err := NewInboundDetourHandlerDynamic(ctx, inbound)
			if err != nil {
				log.Error("V2Ray: Failed to create detour handler: ", err)
				return nil, common.ErrBadConfiguration
			}
			inboundHandler = dh
		default:
			log.Error("V2Ray: Unknown allocation strategy: ", allocConfig.Type)
			return nil, common.ErrBadConfiguration
		}
		vpoint.inboundHandlers = append(vpoint.inboundHandlers, inboundHandler)
		if len(inbound.Tag) > 0 {
			vpoint.taggedInboundHandlers[inbound.Tag] = inboundHandler
		}
	}

	vpoint.outboundHandlers = make([]proxy.OutboundHandler, 0, 8)
	vpoint.taggedOutboundHandlers = make(map[string]proxy.OutboundHandler)
	for idx, outbound := range pConfig.Outbound {
		outboundSettings, err := outbound.GetTypedSettings()
		if err != nil {
			return nil, err
		}
		outboundHandler, err := proxy.CreateOutboundHandler(proxy.ContextWithOutboundMeta(ctx, &proxy.OutboundHandlerMeta{
			Tag:            outbound.Tag,
			Address:        outbound.GetSendThroughValue(),
			StreamSettings: outbound.StreamSettings,
			ProxySettings:  outbound.ProxySettings,
		}), outboundSettings)
		if err != nil {
			log.Error("V2Ray: Failed to create detour outbound connection handler: ", err)
			return nil, err
		}
		if idx == 0 {
			outboundHandlerManager.SetDefaultHandler(outboundHandler)
		}
		if len(outbound.Tag) > 0 {
			outboundHandlerManager.SetHandler(outbound.Tag, outboundHandler)
			vpoint.taggedOutboundHandlers[outbound.Tag] = outboundHandler
		}

		vpoint.outboundHandlers = append(vpoint.outboundHandlers, outboundHandler)
	}

	if err := vpoint.space.Initialize(); err != nil {
		return nil, err
	}

	return vpoint, nil
}

func (Point) Interface() interface{} {
	return (*proxyman.InboundHandlerManager)(nil)
}

func (v *Point) Close() {
	for _, inbound := range v.inboundHandlers {
		inbound.Close()
	}
}

// Start starts the Point server, and return any error during the process.
// In the case of any errors, the state of the server is unpredicatable.
func (v *Point) Start() error {
	for _, inbound := range v.inboundHandlers {
		err := inbound.Start()
		if err != nil {
			return err
		}
	}
	log.Warning("V2Ray started.")

	return nil
}

func (v *Point) GetHandler(tag string) (proxy.InboundHandler, int) {
	handler, found := v.taggedInboundHandlers[tag]
	if !found {
		log.Warning("V2Ray: Unable to find an inbound handler with tag: ", tag)
		return nil, 0
	}
	return handler.GetConnectionHandler()
}

func (v *Point) Release() {

}
