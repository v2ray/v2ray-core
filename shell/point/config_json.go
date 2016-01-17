// +build json

package point

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"

	"github.com/v2ray/v2ray-core/app/router"
	"github.com/v2ray/v2ray-core/common/log"
	v2net "github.com/v2ray/v2ray-core/common/net"
)

func (this *Config) UnmarshalJSON(data []byte) error {
	type JsonConfig struct {
		Port            v2net.Port              `json:"port"` // Port of this Point server.
		LogConfig       *LogConfig              `json:"log"`
		RouterConfig    *router.Config          `json:"routing"`
		InboundConfig   *ConnectionConfig       `json:"inbound"`
		OutboundConfig  *ConnectionConfig       `json:"outbound"`
		InboundDetours  []*InboundDetourConfig  `json:"inboundDetour"`
		OutboundDetours []*OutboundDetourConfig `json:"outboundDetour"`
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
	return nil
}

func (this *ConnectionConfig) UnmarshalJSON(data []byte) error {
	type JsonConnectionConfig struct {
		Protocol string          `json:"protocol"`
		Settings json.RawMessage `json:"settings"`
	}
	jsonConfig := new(JsonConnectionConfig)
	if err := json.Unmarshal(data, jsonConfig); err != nil {
		return err
	}
	this.Protocol = jsonConfig.Protocol
	this.Settings = jsonConfig.Settings
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
	default:
		this.LogLevel = log.WarningLevel
	}
	return nil
}

func (this *InboundDetourAllocationConfig) UnmarshalJSON(data []byte) error {
	type JsonInboundDetourAllocationConfig struct {
		Strategy    string `json:"strategy"`
		Concurrency int    `json:"concurrency"`
		RefreshSec  int    `json:"refresh"`
	}
	jsonConfig := new(JsonInboundDetourAllocationConfig)
	if err := json.Unmarshal(data, jsonConfig); err != nil {
		return err
	}
	this.Strategy = jsonConfig.Strategy
	this.Concurrency = jsonConfig.Concurrency
	this.Refresh = jsonConfig.RefreshSec
	return nil
}

func (this *InboundDetourConfig) UnmarshalJSON(data []byte) error {
	type JsonInboundDetourConfig struct {
		Protocol   string                         `json:"protocol"`
		PortRange  *v2net.PortRange               `json:"port"`
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
		return BadConfiguration
	}
	this.Protocol = jsonConfig.Protocol
	this.PortRange = *jsonConfig.PortRange
	this.Settings = jsonConfig.Settings
	this.Tag = jsonConfig.Tag
	this.Allocation = jsonConfig.Allocation
	return nil
}

func (this *OutboundDetourConfig) UnmarshalJSON(data []byte) error {
	type JsonOutboundDetourConfig struct {
		Protocol string          `json:"protocol"`
		Tag      string          `json:"tag"`
		Settings json.RawMessage `json:"settings"`
	}
	jsonConfig := new(JsonOutboundDetourConfig)
	if err := json.Unmarshal(data, jsonConfig); err != nil {
		return err
	}
	this.Protocol = jsonConfig.Protocol
	this.Tag = jsonConfig.Tag
	this.Settings = jsonConfig.Settings
	return nil
}

func JsonLoadConfig(file string) (*Config, error) {
	fixedFile := os.ExpandEnv(file)
	rawConfig, err := ioutil.ReadFile(fixedFile)
	if err != nil {
		log.Error("Failed to read server config file (%s): %v", file, err)
		return nil, err
	}

	jsonConfig := &Config{}
	err = json.Unmarshal(rawConfig, jsonConfig)
	if err != nil {
		log.Error("Failed to load server config: %v", err)
		return nil, err
	}

	return jsonConfig, err
}

func init() {
	configLoader = JsonLoadConfig
}
