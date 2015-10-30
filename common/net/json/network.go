package json

import (
	"encoding/json"
	"strings"

	v2net "github.com/v2ray/v2ray-core/common/net"
)

type NetworkList []string

func (this *NetworkList) UnmarshalJSON(data []byte) error {
	var strList []string
	err := json.Unmarshal(data, &strList)
	if err != nil {
		return err
	}
	*this = make([]string, len(strList))
	for idx, str := range strList {
		(*this)[idx] = strings.ToLower(str)
	}
	return nil
}

func (this *NetworkList) HasNetwork(network v2net.Network) bool {
	for _, value := range *this {
		if value == string(network) {
			return true
		}
	}
	return false
}
