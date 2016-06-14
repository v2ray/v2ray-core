// +build json

package internet

import (
	"encoding/json"

	v2net "github.com/v2ray/v2ray-core/common/net"
)

func (this *StreamSettings) UnmarshalJSON(data []byte) error {
	type JSONConfig struct {
		Network v2net.NetworkList `json:"network"`
	}
	this.Type = StreamConnectionTypeRawTCP
	jsonConfig := new(JSONConfig)
	if err := json.Unmarshal(data, jsonConfig); err != nil {
		return err
	}
	if jsonConfig.Network.HasNetwork(v2net.KCPNetwork) {
		this.Type |= StreamConnectionTypeKCP
	}
	if jsonConfig.Network.HasNetwork(v2net.TCPNetwork) {
		this.Type |= StreamConnectionTypeTCP
	}
	return nil
}
