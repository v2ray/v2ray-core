package conf

import (
	"encoding/json"
	"errors"
	"strings"

	"io"
	"v2ray.com/core"
	"v2ray.com/core/common/loader"
	v2net "v2ray.com/core/common/net"
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
	}, "protocol", "settings")
)

type InboundConnectionConfig struct {
	Port          uint16          `json:"port"`
	Listen        *Address        `json:"listen"`
	Protocol      string          `json:"protocol"`
	StreamSetting *StreamConfig   `json:"streamSettings"`
	Settings      json.RawMessage `json:"settings"`
	AllowPassive  bool            `json:"allowPassive"`
}

func (this *InboundConnectionConfig) Build() (*core.InboundConnectionConfig, error) {
	config := new(core.InboundConnectionConfig)
	config.PortRange = &v2net.PortRange{
		From: uint32(this.Port),
		To:   uint32(this.Port),
	}
	if this.Listen != nil {
		if this.Listen.Family().IsDomain() {
			return nil, errors.New("Point: Unable to listen on domain address: " + this.Listen.Domain())
		}
		config.ListenOn = this.Listen.Build()
	}
	if this.StreamSetting != nil {
		ts, err := this.StreamSetting.Build()
		if err != nil {
			return nil, err
		}
		config.StreamSettings = ts
	}
	config.AllowPassiveConnection = this.AllowPassive

	jsonConfig, err := inboundConfigLoader.LoadWithID(this.Settings, this.Protocol)
	if err != nil {
		return nil, errors.New("Failed to load inbound config: " + err.Error())
	}
	ts, err := jsonConfig.(Buildable).Build()
	if err != nil {
		return nil, err
	}
	config.Settings = ts
	return config, nil
}

type OutboundConnectionConfig struct {
	Protocol      string          `json:"protocol"`
	SendThrough   *Address        `json:"sendThrough"`
	StreamSetting *StreamConfig   `json:"streamSettings"`
	ProxySettings *ProxyConfig    `json:"proxySettings"`
	Settings      json.RawMessage `json:"settings"`
}

func (this *OutboundConnectionConfig) Build() (*core.OutboundConnectionConfig, error) {
	config := new(core.OutboundConnectionConfig)
	rawConfig, err := outboundConfigLoader.LoadWithID(this.Settings, this.Protocol)
	if err != nil {
		return nil, errors.New("Failed to parse outbound config: " + err.Error())
	}
	ts, err := rawConfig.(Buildable).Build()
	if err != nil {
		return nil, err
	}
	config.Settings = ts

	if this.SendThrough != nil {
		address := this.SendThrough
		if address.Family().IsDomain() {
			return nil, errors.New("Point: Unable to send through: " + address.String())
		}
		config.SendThrough = address.Build()
	}
	if this.StreamSetting != nil {
		ss, err := this.StreamSetting.Build()
		if err != nil {
			return nil, err
		}
		config.StreamSettings = ss
	}
	if this.ProxySettings != nil {
		ps, err := this.ProxySettings.Build()
		if err != nil {
			return nil, errors.New("Outbound: invalid proxy settings: " + err.Error())
		}
		config.ProxySettings = ps
	}
	return config, nil
}

type InboundDetourAllocationConfig struct {
	Strategy    string  `json:"strategy"`
	Concurrency *uint32 `json:"concurrency"`
	RefreshMin  *uint32 `json:"refresh"`
}

func (this *InboundDetourAllocationConfig) Build() (*core.AllocationStrategy, error) {
	config := new(core.AllocationStrategy)
	switch strings.ToLower(this.Strategy) {
	case "always":
		config.Type = core.AllocationStrategy_Always
	case "random":
		config.Type = core.AllocationStrategy_Random
	case "external":
		config.Type = core.AllocationStrategy_External
	default:
		return nil, errors.New("Unknown allocation strategy: " + this.Strategy)
	}
	if this.Concurrency != nil {
		config.Concurrency = &core.AllocationStrategyConcurrency{
			Value: *this.Concurrency,
		}
	}

	if this.RefreshMin != nil {
		config.Refresh = &core.AllocationStrategyRefresh{
			Value: *this.RefreshMin,
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

func (this *InboundDetourConfig) Build() (*core.InboundConnectionConfig, error) {
	config := new(core.InboundConnectionConfig)
	if this.PortRange == nil {
		return nil, errors.New("Point: Port range not specified in InboundDetour.")
	}
	config.PortRange = this.PortRange.Build()

	if this.ListenOn != nil {
		if this.ListenOn.Family().IsDomain() {
			return nil, errors.New("Point: Unable to listen on domain address: " + this.ListenOn.Domain())
		}
		config.ListenOn = this.ListenOn.Build()
	}
	config.Tag = this.Tag
	if this.Allocation != nil {
		as, err := this.Allocation.Build()
		if err != nil {
			return nil, err
		}
		config.AllocationStrategy = as
	}
	if this.StreamSetting != nil {
		ss, err := this.StreamSetting.Build()
		if err != nil {
			return nil, err
		}
		config.StreamSettings = ss
	}
	config.AllowPassiveConnection = this.AllowPassive

	rawConfig, err := inboundConfigLoader.LoadWithID(this.Settings, this.Protocol)
	if err != nil {
		return nil, errors.New("Failed to load inbound detour config: " + err.Error())
	}
	ts, err := rawConfig.(Buildable).Build()
	if err != nil {
		return nil, err
	}
	config.Settings = ts
	return config, nil
}

type OutboundDetourConfig struct {
	Protocol      string          `json:"protocol"`
	SendThrough   *Address        `json:"sendThrough"`
	Tag           string          `json:"tag"`
	Settings      json.RawMessage `json:"settings"`
	StreamSetting *StreamConfig   `json:"streamSettings"`
	ProxySettings *ProxyConfig    `json:"proxySettings"`
}

func (this *OutboundDetourConfig) Build() (*core.OutboundConnectionConfig, error) {
	config := new(core.OutboundConnectionConfig)
	config.Tag = this.Tag

	if this.SendThrough != nil {
		address := this.SendThrough
		if address.Family().IsDomain() {
			return nil, errors.New("Point: Unable to send through: " + address.String())
		}
		config.SendThrough = address.Build()
	}

	if this.StreamSetting != nil {
		ss, err := this.StreamSetting.Build()
		if err != nil {
			return nil, err
		}
		config.StreamSettings = ss
	}

	rawConfig, err := outboundConfigLoader.LoadWithID(this.Settings, this.Protocol)
	if err != nil {
		return nil, errors.New("Failed to parse to outbound detour config: " + err.Error())
	}
	ts, err := rawConfig.(Buildable).Build()
	if err != nil {
		return nil, err
	}

	if this.ProxySettings != nil {
		ps, err := this.ProxySettings.Build()
		if err != nil {
			return nil, errors.New("OutboundDetour: invalid proxy settings: " + err.Error())
		}
		config.ProxySettings = ps
	}
	config.Settings = ts
	return config, nil
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

func (this *Config) Build() (*core.Config, error) {
	config := new(core.Config)

	if this.LogConfig != nil {
		config.Log = this.LogConfig.Build()
	}

	if this.Transport != nil {
		ts, err := this.Transport.Build()
		if err != nil {
			return nil, err
		}
		config.Transport = ts
	}

	if this.RouterConfig != nil {
		routerConfig, err := this.RouterConfig.Build()
		if err != nil {
			return nil, err
		}
		config.App = append(config.App, loader.NewTypedSettings(routerConfig))
	}

	if this.DNSConfig != nil {
		config.App = append(config.App, loader.NewTypedSettings(this.DNSConfig.Build()))
	}

	if this.InboundConfig == nil {
		return nil, errors.New("No inbound config specified.")
	}

	if this.InboundConfig.Port == 0 && this.Port > 0 {
		this.InboundConfig.Port = this.Port
	}

	ic, err := this.InboundConfig.Build()
	if err != nil {
		return nil, err
	}
	config.Inbound = append(config.Inbound, ic)

	for _, rawInboundConfig := range this.InboundDetours {
		ic, err := rawInboundConfig.Build()
		if err != nil {
			return nil, err
		}
		config.Inbound = append(config.Inbound, ic)
	}

	oc, err := this.OutboundConfig.Build()
	if err != nil {
		return nil, err
	}
	config.Outbound = append(config.Outbound, oc)

	for _, rawOutboundConfig := range this.OutboundDetours {
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
		decoder := json.NewDecoder(input)
		err := decoder.Decode(jsonConfig)
		if err != nil {
			return nil, errors.New("Point: Failed to load server config: " + err.Error())
		}

		return jsonConfig.Build()
	})
}
