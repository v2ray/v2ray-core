package conf

import (
	"encoding/json"
	"io"
	"strings"

	"v2ray.com/core"
	"v2ray.com/core/app/proxyman"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/common/serial"
	json_reader "v2ray.com/core/tools/conf/json"
)

var (
	inboundConfigLoader = NewJSONConfigLoader(ConfigCreatorCache{
		"dokodemo-door": func() interface{} { return new(DokodemoConfig) },
		"http":          func() interface{} { return new(HttpServerConfig) },
		"shadowsocks":   func() interface{} { return new(ShadowsocksServerConfig) },
		"socks":         func() interface{} { return new(SocksServerConfig) },
		"vmess":         func() interface{} { return new(VMessInboundConfig) },
	}, "protocol", "settings")

	outboundConfigLoader = NewJSONConfigLoader(ConfigCreatorCache{
		"blackhole":   func() interface{} { return new(BlackholeConfig) },
		"freedom":     func() interface{} { return new(FreedomConfig) },
		"shadowsocks": func() interface{} { return new(ShadowsocksClientConfig) },
		"vmess":       func() interface{} { return new(VMessOutboundConfig) },
		"socks":       func() interface{} { return new(SocksClientConfig) },
	}, "protocol", "settings")
)

func toProtocolList(s []string) ([]proxyman.KnownProtocols, error) {
	kp := make([]proxyman.KnownProtocols, 0, 8)
	for _, p := range s {
		switch strings.ToLower(p) {
		case "http":
			kp = append(kp, proxyman.KnownProtocols_HTTP)
		case "https", "tls", "ssl":
			kp = append(kp, proxyman.KnownProtocols_TLS)
		default:
			return nil, newError("Unknown protocol: ", p)
		}
	}
	return kp, nil
}

type InboundConnectionConfig struct {
	Port           uint16          `json:"port"`
	Listen         *Address        `json:"listen"`
	Protocol       string          `json:"protocol"`
	StreamSetting  *StreamConfig   `json:"streamSettings"`
	Settings       json.RawMessage `json:"settings"`
	Tag            string          `json:"tag"`
	DomainOverride *StringList     `json:"domainOverride"`
}

func (v *InboundConnectionConfig) Build() (*proxyman.InboundHandlerConfig, error) {
	receiverConfig := &proxyman.ReceiverConfig{
		PortRange: &v2net.PortRange{
			From: uint32(v.Port),
			To:   uint32(v.Port),
		},
	}
	if v.Listen != nil {
		if v.Listen.Family().IsDomain() {
			return nil, newError("unable to listen on domain address: " + v.Listen.Domain())
		}
		receiverConfig.Listen = v.Listen.Build()
	}
	if v.StreamSetting != nil {
		ts, err := v.StreamSetting.Build()
		if err != nil {
			return nil, err
		}
		receiverConfig.StreamSettings = ts
	}
	if v.DomainOverride != nil {
		kp, err := toProtocolList(*v.DomainOverride)
		if err != nil {
			return nil, newError("failed to parse inbound config").Base(err)
		}
		receiverConfig.DomainOverride = kp
	}

	jsonConfig, err := inboundConfigLoader.LoadWithID(v.Settings, v.Protocol)
	if err != nil {
		return nil, newError("failed to load inbound config.").Base(err)
	}
	if dokodemoConfig, ok := jsonConfig.(*DokodemoConfig); ok {
		receiverConfig.ReceiveOriginalDestination = dokodemoConfig.Redirect
	}
	ts, err := jsonConfig.(Buildable).Build()
	if err != nil {
		return nil, err
	}

	return &proxyman.InboundHandlerConfig{
		Tag:              v.Tag,
		ReceiverSettings: serial.ToTypedMessage(receiverConfig),
		ProxySettings:    ts,
	}, nil
}

type MuxConfig struct {
	Enabled     bool   `json:"enabled"`
	Concurrency uint16 `json:"concurrency"`
}

func (c *MuxConfig) GetConcurrency() uint16 {
	if c.Concurrency == 0 {
		return 8
	}
	return c.Concurrency
}

type OutboundConnectionConfig struct {
	Protocol      string          `json:"protocol"`
	SendThrough   *Address        `json:"sendThrough"`
	StreamSetting *StreamConfig   `json:"streamSettings"`
	ProxySettings *ProxyConfig    `json:"proxySettings"`
	Settings      json.RawMessage `json:"settings"`
	Tag           string          `json:"tag"`
	MuxSettings   *MuxConfig      `json:"mux"`
}

func (v *OutboundConnectionConfig) Build() (*proxyman.OutboundHandlerConfig, error) {
	senderSettings := &proxyman.SenderConfig{}

	if v.SendThrough != nil {
		address := v.SendThrough
		if address.Family().IsDomain() {
			return nil, newError("invalid sendThrough address: " + address.String())
		}
		senderSettings.Via = address.Build()
	}
	if v.StreamSetting != nil {
		ss, err := v.StreamSetting.Build()
		if err != nil {
			return nil, err
		}
		senderSettings.StreamSettings = ss
	}
	if v.ProxySettings != nil {
		ps, err := v.ProxySettings.Build()
		if err != nil {
			return nil, newError("invalid outbound proxy settings").Base(err)
		}
		senderSettings.ProxySettings = ps
	}

	if v.MuxSettings != nil && v.MuxSettings.Enabled {
		senderSettings.MultiplexSettings = &proxyman.MultiplexingConfig{
			Enabled:     true,
			Concurrency: uint32(v.MuxSettings.GetConcurrency()),
		}
	}

	rawConfig, err := outboundConfigLoader.LoadWithID(v.Settings, v.Protocol)
	if err != nil {
		return nil, newError("failed to parse outbound config").Base(err)
	}
	ts, err := rawConfig.(Buildable).Build()
	if err != nil {
		return nil, err
	}

	return &proxyman.OutboundHandlerConfig{
		SenderSettings: serial.ToTypedMessage(senderSettings),
		ProxySettings:  ts,
		Tag:            v.Tag,
	}, nil
}

type InboundDetourAllocationConfig struct {
	Strategy    string  `json:"strategy"`
	Concurrency *uint32 `json:"concurrency"`
	RefreshMin  *uint32 `json:"refresh"`
}

func (v *InboundDetourAllocationConfig) Build() (*proxyman.AllocationStrategy, error) {
	config := new(proxyman.AllocationStrategy)
	switch strings.ToLower(v.Strategy) {
	case "always":
		config.Type = proxyman.AllocationStrategy_Always
	case "random":
		config.Type = proxyman.AllocationStrategy_Random
	case "external":
		config.Type = proxyman.AllocationStrategy_External
	default:
		return nil, newError("unknown allocation strategy: ", v.Strategy)
	}
	if v.Concurrency != nil {
		config.Concurrency = &proxyman.AllocationStrategy_AllocationStrategyConcurrency{
			Value: *v.Concurrency,
		}
	}

	if v.RefreshMin != nil {
		config.Refresh = &proxyman.AllocationStrategy_AllocationStrategyRefresh{
			Value: *v.RefreshMin,
		}
	}

	return config, nil
}

type InboundDetourConfig struct {
	Protocol       string                         `json:"protocol"`
	PortRange      *PortRange                     `json:"port"`
	ListenOn       *Address                       `json:"listen"`
	Settings       json.RawMessage                `json:"settings"`
	Tag            string                         `json:"tag"`
	Allocation     *InboundDetourAllocationConfig `json:"allocate"`
	StreamSetting  *StreamConfig                  `json:"streamSettings"`
	DomainOverride *StringList                    `json:"domainOverride"`
}

func (v *InboundDetourConfig) Build() (*proxyman.InboundHandlerConfig, error) {
	receiverSettings := &proxyman.ReceiverConfig{}

	if v.PortRange == nil {
		return nil, newError("port range not specified in InboundDetour.")
	}
	receiverSettings.PortRange = v.PortRange.Build()

	if v.ListenOn != nil {
		if v.ListenOn.Family().IsDomain() {
			return nil, newError("unable to listen on domain address: ", v.ListenOn.Domain())
		}
		receiverSettings.Listen = v.ListenOn.Build()
	}
	if v.Allocation != nil {
		as, err := v.Allocation.Build()
		if err != nil {
			return nil, err
		}
		receiverSettings.AllocationStrategy = as
	}
	if v.StreamSetting != nil {
		ss, err := v.StreamSetting.Build()
		if err != nil {
			return nil, err
		}
		receiverSettings.StreamSettings = ss
	}
	if v.DomainOverride != nil {
		kp, err := toProtocolList(*v.DomainOverride)
		if err != nil {
			return nil, newError("failed to parse inbound detour config").Base(err)
		}
		receiverSettings.DomainOverride = kp
	}

	rawConfig, err := inboundConfigLoader.LoadWithID(v.Settings, v.Protocol)
	if err != nil {
		return nil, newError("failed to load inbound detour config.").Base(err)
	}
	if dokodemoConfig, ok := rawConfig.(*DokodemoConfig); ok {
		receiverSettings.ReceiveOriginalDestination = dokodemoConfig.Redirect
	}
	ts, err := rawConfig.(Buildable).Build()
	if err != nil {
		return nil, err
	}

	return &proxyman.InboundHandlerConfig{
		Tag:              v.Tag,
		ReceiverSettings: serial.ToTypedMessage(receiverSettings),
		ProxySettings:    ts,
	}, nil
}

type OutboundDetourConfig struct {
	Protocol      string          `json:"protocol"`
	SendThrough   *Address        `json:"sendThrough"`
	Tag           string          `json:"tag"`
	Settings      json.RawMessage `json:"settings"`
	StreamSetting *StreamConfig   `json:"streamSettings"`
	ProxySettings *ProxyConfig    `json:"proxySettings"`
	MuxSettings   *MuxConfig      `json:"mux"`
}

func (v *OutboundDetourConfig) Build() (*proxyman.OutboundHandlerConfig, error) {
	senderSettings := &proxyman.SenderConfig{}

	if v.SendThrough != nil {
		address := v.SendThrough
		if address.Family().IsDomain() {
			return nil, newError("unable to send through: " + address.String())
		}
		senderSettings.Via = address.Build()
	}

	if v.StreamSetting != nil {
		ss, err := v.StreamSetting.Build()
		if err != nil {
			return nil, err
		}
		senderSettings.StreamSettings = ss
	}

	if v.ProxySettings != nil {
		ps, err := v.ProxySettings.Build()
		if err != nil {
			return nil, newError("invalid outbound detour proxy settings.").Base(err)
		}
		senderSettings.ProxySettings = ps
	}

	if v.MuxSettings != nil && v.MuxSettings.Enabled {
		senderSettings.MultiplexSettings = &proxyman.MultiplexingConfig{
			Enabled:     true,
			Concurrency: uint32(v.MuxSettings.GetConcurrency()),
		}
	}

	rawConfig, err := outboundConfigLoader.LoadWithID(v.Settings, v.Protocol)
	if err != nil {
		return nil, newError("failed to parse to outbound detour config.").Base(err)
	}
	ts, err := rawConfig.(Buildable).Build()
	if err != nil {
		return nil, err
	}

	return &proxyman.OutboundHandlerConfig{
		SenderSettings: serial.ToTypedMessage(senderSettings),
		Tag:            v.Tag,
		ProxySettings:  ts,
	}, nil
}

type Config struct {
	Port            uint16                    `json:"port"` // Port of this Point server.
	LogConfig       *LogConfig                `json:"log"`
	RouterConfig    *RouterConfig             `json:"routing"`
	DNSConfig       *DnsConfig                `json:"dns"`
	InboundConfig   *InboundConnectionConfig  `json:"inbound"`
	OutboundConfig  *OutboundConnectionConfig `json:"outbound"`
	InboundDetours  []InboundDetourConfig     `json:"inboundDetour"`
	OutboundDetours []OutboundDetourConfig    `json:"outboundDetour"`
	Transport       *TransportConfig          `json:"transport"`
}

func (v *Config) Build() (*core.Config, error) {
	config := new(core.Config)

	if v.LogConfig != nil {
		config.App = append(config.App, serial.ToTypedMessage(v.LogConfig.Build()))
	}

	if v.Transport != nil {
		ts, err := v.Transport.Build()
		if err != nil {
			return nil, err
		}
		config.Transport = ts
	}

	if v.RouterConfig != nil {
		routerConfig, err := v.RouterConfig.Build()
		if err != nil {
			return nil, err
		}
		config.App = append(config.App, serial.ToTypedMessage(routerConfig))
	}

	if v.DNSConfig != nil {
		config.App = append(config.App, serial.ToTypedMessage(v.DNSConfig.Build()))
	}

	if v.InboundConfig == nil {
		return nil, newError("no inbound config specified")
	}

	if v.InboundConfig.Port == 0 && v.Port > 0 {
		v.InboundConfig.Port = v.Port
	}

	ic, err := v.InboundConfig.Build()
	if err != nil {
		return nil, err
	}
	config.Inbound = append(config.Inbound, ic)

	for _, rawInboundConfig := range v.InboundDetours {
		ic, err := rawInboundConfig.Build()
		if err != nil {
			return nil, err
		}
		config.Inbound = append(config.Inbound, ic)
	}

	if v.OutboundConfig == nil {
		return nil, newError("no outbound config specified")
	}
	oc, err := v.OutboundConfig.Build()
	if err != nil {
		return nil, err
	}
	config.Outbound = append(config.Outbound, oc)

	for _, rawOutboundConfig := range v.OutboundDetours {
		oc, err := rawOutboundConfig.Build()
		if err != nil {
			return nil, err
		}
		config.Outbound = append(config.Outbound, oc)
	}

	return config, nil
}

func init() {
	core.RegisterConfigLoader(core.ConfigFormat_JSON, func(input io.Reader) (*core.Config, error) {
		jsonConfig := &Config{}
		decoder := json.NewDecoder(&json_reader.Reader{
			Reader: input,
		})
		err := decoder.Decode(jsonConfig)
		if err != nil {
			return nil, newError("invalid V2Ray config").Base(err)
		}

		return jsonConfig.Build()
	})
}
