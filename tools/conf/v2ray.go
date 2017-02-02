package conf

import (
	"encoding/json"
	"io"
	"strings"

	"v2ray.com/core"
	"v2ray.com/core/app/proxyman"
	"v2ray.com/core/common/errors"
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

type InboundConnectionConfig struct {
	Port          uint16          `json:"port"`
	Listen        *Address        `json:"listen"`
	Protocol      string          `json:"protocol"`
	StreamSetting *StreamConfig   `json:"streamSettings"`
	Settings      json.RawMessage `json:"settings"`
	AllowPassive  bool            `json:"allowPassive"`
	Tag           string          `json:"tag"`
}


func ImportJsonParser(){
}

func (v *InboundConnectionConfig) Build() (*proxyman.InboundHandlerConfig, error) {
	receiverConfig := &proxyman.ReceiverConfig{
		PortRange: &v2net.PortRange{
			From: uint32(v.Port),
			To:   uint32(v.Port),
		},
		AllowPassiveConnection: v.AllowPassive,
	}
	if v.Listen != nil {
		if v.Listen.Family().IsDomain() {
			return nil, errors.New("Config: Unable to listen on domain address: " + v.Listen.Domain())
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

	jsonConfig, err := inboundConfigLoader.LoadWithID(v.Settings, v.Protocol)
	if err != nil {
		return nil, errors.Base(err).Message("Config: Failed to load inbound config.")
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

type OutboundConnectionConfig struct {
	Protocol      string          `json:"protocol"`
	SendThrough   *Address        `json:"sendThrough"`
	StreamSetting *StreamConfig   `json:"streamSettings"`
	ProxySettings *ProxyConfig    `json:"proxySettings"`
	Settings      json.RawMessage `json:"settings"`
	Tag           string          `json:"tag"`
}

func (v *OutboundConnectionConfig) Build() (*proxyman.OutboundHandlerConfig, error) {
	senderSettings := &proxyman.SenderConfig{}

	if v.SendThrough != nil {
		address := v.SendThrough
		if address.Family().IsDomain() {
			return nil, errors.New("Config: Invalid sendThrough address: " + address.String())
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
			return nil, errors.Base(err).Message("Config: Invalid outbound proxy settings.")
		}
		senderSettings.ProxySettings = ps
	}

	rawConfig, err := outboundConfigLoader.LoadWithID(v.Settings, v.Protocol)
	if err != nil {
		return nil, errors.Base(err).Message("Config: Failed to parse outbound config.")
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
		return nil, errors.New("Config: Unknown allocation strategy: ", v.Strategy)
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
	Protocol      string                         `json:"protocol"`
	PortRange     *PortRange                     `json:"port"`
	ListenOn      *Address                       `json:"listen"`
	Settings      json.RawMessage                `json:"settings"`
	Tag           string                         `json:"tag"`
	Allocation    *InboundDetourAllocationConfig `json:"allocate"`
	StreamSetting *StreamConfig                  `json:"streamSettings"`
	AllowPassive  bool                           `json:"allowPassive"`
}

func (v *InboundDetourConfig) Build() (*proxyman.InboundHandlerConfig, error) {
	receiverSettings := &proxyman.ReceiverConfig{
		AllowPassiveConnection: v.AllowPassive,
	}

	if v.PortRange == nil {
		return nil, errors.New("Config: Port range not specified in InboundDetour.")
	}
	receiverSettings.PortRange = v.PortRange.Build()

	if v.ListenOn != nil {
		if v.ListenOn.Family().IsDomain() {
			return nil, errors.New("Config: Unable to listen on domain address: ", v.ListenOn.Domain())
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

	rawConfig, err := inboundConfigLoader.LoadWithID(v.Settings, v.Protocol)
	if err != nil {
		return nil, errors.Base(err).Message("Config: Failed to load inbound detour config.")
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
}

func (v *OutboundDetourConfig) Build() (*proxyman.OutboundHandlerConfig, error) {
	senderSettings := &proxyman.SenderConfig{}

	if v.SendThrough != nil {
		address := v.SendThrough
		if address.Family().IsDomain() {
			return nil, errors.New("Config: Unable to send through: " + address.String())
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
			return nil, errors.Base(err).Message("Config: Invalid outbound detour proxy settings.")
		}
		senderSettings.ProxySettings = ps
	}

	rawConfig, err := outboundConfigLoader.LoadWithID(v.Settings, v.Protocol)
	if err != nil {
		return nil, errors.Base(err).Message("Config: Failed to parse to outbound detour config.")
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
		return nil, errors.New("Config: No inbound config specified.")
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
		return nil, errors.New("Config: No outbound config specified.")
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
			return nil, errors.Base(err).Message("Config: Invalid V2Ray config.")
		}

		return jsonConfig.Build()
	})
}
