package json

import (
	"encoding/json"

	v2net "github.com/v2ray/v2ray-core/common/net"
	v2netjson "github.com/v2ray/v2ray-core/common/net/json"
	"github.com/v2ray/v2ray-core/shell/point"
)

type InboundDetourAllocationConfig struct {
	StrategyValue    string `json:"strategy"`
	ConcurrencyValue int    `json:"concurrency"`
	RefreshSec       int    `json:"refresh"`
}

func (this *InboundDetourAllocationConfig) Refresh() int {
	return this.RefreshSec
}

func (this *InboundDetourAllocationConfig) Strategy() string {
	return this.StrategyValue
}

func (this *InboundDetourAllocationConfig) Concurrency() int {
	return this.ConcurrencyValue
}

type InboundDetourConfig struct {
	ProtocolValue   string                         `json:"protocol"`
	PortRangeValue  *v2netjson.PortRange           `json:"port"`
	SettingsValue   json.RawMessage                `json:"settings"`
	TagValue        string                         `json:"tag"`
	AllocationValue *InboundDetourAllocationConfig `json:"allocate"`
}

func (this *InboundDetourConfig) Allocation() point.InboundDetourAllocationConfig {
	return this.AllocationValue
}

func (this *InboundDetourConfig) Protocol() string {
	return this.ProtocolValue
}

func (this *InboundDetourConfig) PortRange() v2net.PortRange {
	return this.PortRangeValue
}

func (this *InboundDetourConfig) Settings() []byte {
	return []byte(this.SettingsValue)
}

func (this *InboundDetourConfig) Tag() string {
	return this.TagValue
}
