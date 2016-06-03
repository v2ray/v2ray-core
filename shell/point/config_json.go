// +build json

package point

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"strings"

	"github.com/v2ray/v2ray-core/app/dns"
	"github.com/v2ray/v2ray-core/app/router"
	"github.com/v2ray/v2ray-core/common/log"
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/transport"
)

const (
	DefaultRefreshMinute = int(9999)
)

func (this *Config) UnmarshalJSON(data []byte) error {
	type JsonConfig struct {
		Port            v2net.Port                `json:"port"` // Port of this Point server.
		LogConfig       *LogConfig                `json:"log"`
		RouterConfig    *router.Config            `json:"routing"`
		DNSConfig       *dns.Config               `json:"dns"`
		InboundConfig   *InboundConnectionConfig  `json:"inbound"`
		OutboundConfig  *OutboundConnectionConfig `json:"outbound"`
		InboundDetours  []*InboundDetourConfig    `json:"inboundDetour"`
		OutboundDetours []*OutboundDetourConfig   `json:"outboundDetour"`
		Transport       *transport.Config         `json:"transport"`
	}
	jsonConfig := new(JsonConfig)
	if err := json.Unmarshal(data, jsonConfig); err != nil {
		return err
	}
	this.Port = jsonConfig.Port
	this.LogConfig = jsonConfig.LogConfig
	this.RouterConfig = jsonConfig.RouterConfig
	this.InboundConfig = jsonConfig.InboundConfig
	this.OutboundConfig = jsonConfig.OutboundConfig
	this.InboundDetours = jsonConfig.InboundDetours
	this.OutboundDetours = jsonConfig.OutboundDetours
	if jsonConfig.DNSConfig == nil {
		jsonConfig.DNSConfig = &dns.Config{
			NameServers: []v2net.Destination{
				v2net.UDPDestination(v2net.DomainAddress("localhost"), v2net.Port(53)),
			},
		}
	}
	this.DNSConfig = jsonConfig.DNSConfig
	this.TransportConfig = jsonConfig.Transport
	return nil
}

func (this *InboundConnectionConfig) UnmarshalJSON(data []byte) error {
	type JsonConfig struct {
		Port     uint16             `json:"port"`
		Listen   *v2net.AddressJson `json:"listen"`
		Protocol string             `json:"protocol"`
		Settings json.RawMessage    `json:"settings"`
	}

	jsonConfig := new(JsonConfig)
	if err := json.Unmarshal(data, jsonConfig); err != nil {
		return err
	}
	this.Port = v2net.Port(jsonConfig.Port)
	this.ListenOn = v2net.AnyIP
	if jsonConfig.Listen != nil {
		if jsonConfig.Listen.Address.IsDomain() {
			return errors.New("Point: Unable to listen on domain address: " + jsonConfig.Listen.Address.Domain())
		}
		this.ListenOn = jsonConfig.Listen.Address
	}

	this.Protocol = jsonConfig.Protocol
	this.Settings = jsonConfig.Settings
	return nil
}

func (this *OutboundConnectionConfig) UnmarshalJSON(data []byte) error {
	type JsonConnectionConfig struct {
		Protocol    string             `json:"protocol"`
		SendThrough *v2net.AddressJson `json:"sendThrough"`
		Settings    json.RawMessage    `json:"settings"`
	}
	jsonConfig := new(JsonConnectionConfig)
	if err := json.Unmarshal(data, jsonConfig); err != nil {
		return err
	}
	this.Protocol = jsonConfig.Protocol
	this.Settings = jsonConfig.Settings

	if jsonConfig.SendThrough != nil {
		address := jsonConfig.SendThrough.Address
		if address.IsDomain() {
			return errors.New("Point: Unable to send through: " + address.String())
		}
		this.SendThrough = address
	}
	return nil
}

func (this *LogConfig) UnmarshalJSON(data []byte) error {
	type JsonLogConfig struct {
		AccessLog string `json:"access"`
		ErrorLog  string `json:"error"`
		LogLevel  string `json:"loglevel"`
	}
	jsonConfig := new(JsonLogConfig)
	if err := json.Unmarshal(data, jsonConfig); err != nil {
		return err
	}
	this.AccessLog = jsonConfig.AccessLog
	this.ErrorLog = jsonConfig.ErrorLog

	level := strings.ToLower(jsonConfig.LogLevel)
	switch level {
	case "debug":
		this.LogLevel = log.DebugLevel
	case "info":
		this.LogLevel = log.InfoLevel
	case "error":
		this.LogLevel = log.ErrorLevel
	case "none":
		this.LogLevel = log.NoneLevel
	default:
		this.LogLevel = log.WarningLevel
	}
	return nil
}

func (this *InboundDetourAllocationConfig) UnmarshalJSON(data []byte) error {
	type JsonInboundDetourAllocationConfig struct {
		Strategy    string `json:"strategy"`
		Concurrency int    `json:"concurrency"`
		RefreshMin  int    `json:"refresh"`
	}
	jsonConfig := new(JsonInboundDetourAllocationConfig)
	if err := json.Unmarshal(data, jsonConfig); err != nil {
		return err
	}
	this.Strategy = jsonConfig.Strategy
	this.Concurrency = jsonConfig.Concurrency
	this.Refresh = jsonConfig.RefreshMin
	if this.Strategy == AllocationStrategyRandom {
		if this.Refresh == 0 {
			this.Refresh = 5
		}
		if this.Concurrency == 0 {
			this.Concurrency = 3
		}
	}
	if this.Refresh == 0 {
		this.Refresh = DefaultRefreshMinute
	}
	return nil
}

func (this *InboundDetourConfig) UnmarshalJSON(data []byte) error {
	type JsonInboundDetourConfig struct {
		Protocol   string                         `json:"protocol"`
		PortRange  *v2net.PortRange               `json:"port"`
		ListenOn   *v2net.AddressJson             `json:"listen"`
		Settings   json.RawMessage                `json:"settings"`
		Tag        string                         `json:"tag"`
		Allocation *InboundDetourAllocationConfig `json:"allocate"`
	}
	jsonConfig := new(JsonInboundDetourConfig)
	if err := json.Unmarshal(data, jsonConfig); err != nil {
		return err
	}
	if jsonConfig.PortRange == nil {
		log.Error("Point: Port range not specified in InboundDetour.")
		return ErrorBadConfiguration
	}
	this.ListenOn = v2net.AnyIP
	if jsonConfig.ListenOn != nil {
		if jsonConfig.ListenOn.Address.IsDomain() {
			return errors.New("Point: Unable to listen on domain address: " + jsonConfig.ListenOn.Address.Domain())
		}
		this.ListenOn = jsonConfig.ListenOn.Address
	}
	this.Protocol = jsonConfig.Protocol
	this.PortRange = *jsonConfig.PortRange
	this.Settings = jsonConfig.Settings
	this.Tag = jsonConfig.Tag
	this.Allocation = jsonConfig.Allocation
	if this.Allocation == nil {
		this.Allocation = &InboundDetourAllocationConfig{
			Strategy: AllocationStrategyAlways,
			Refresh:  DefaultRefreshMinute,
		}
	}
	return nil
}

func (this *OutboundDetourConfig) UnmarshalJSON(data []byte) error {
	type JsonOutboundDetourConfig struct {
		Protocol    string             `json:"protocol"`
		SendThrough *v2net.AddressJson `json:"sendThrough"`
		Tag         string             `json:"tag"`
		Settings    json.RawMessage    `json:"settings"`
	}
	jsonConfig := new(JsonOutboundDetourConfig)
	if err := json.Unmarshal(data, jsonConfig); err != nil {
		return err
	}
	this.Protocol = jsonConfig.Protocol
	this.Tag = jsonConfig.Tag
	this.Settings = jsonConfig.Settings

	if jsonConfig.SendThrough != nil {
		address := jsonConfig.SendThrough.Address
		if address.IsDomain() {
			return errors.New("Point: Unable to send through: " + address.String())
		}
		this.SendThrough = address
	}
	return nil
}

func JsonLoadConfig(file string) (*Config, error) {
	fixedFile := os.ExpandEnv(file)
	rawConfig, err := ioutil.ReadFile(fixedFile)
	if err != nil {
		log.Error("Failed to read server config file (", file, "): ", file, err)
		return nil, err
	}

	jsonConfig := &Config{}
	err = json.Unmarshal(rawConfig, jsonConfig)
	if err != nil {
		log.Error("Failed to load server config: ", err)
		return nil, err
	}

	return jsonConfig, err
}

func init() {
	configLoader = JsonLoadConfig
}
