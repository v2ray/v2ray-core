// +build json

package outbound

import (
	"encoding/json"
	"errors"

	"github.com/v2ray/v2ray-core/common/log"
	"github.com/v2ray/v2ray-core/proxy/internal"
)

func (this *Config) UnmarshalJSON(data []byte) error {
	type RawOutbound struct {
		Receivers []*Receiver `json:"vnext"`
	}
	rawOutbound := &RawOutbound{}
	err := json.Unmarshal(data, rawOutbound)
	if err != nil {
		return errors.New("VMessOut: Failed to parse config: " + err.Error())
	}
	if len(rawOutbound.Receivers) == 0 {
		log.Error("VMessOut: 0 VMess receiver configured.")
		return internal.ErrorBadConfiguration
	}
	this.Receivers = rawOutbound.Receivers
	return nil
}

func init() {
	internal.RegisterOutboundConfig("vmess", func() interface{} { return new(Config) })
}
