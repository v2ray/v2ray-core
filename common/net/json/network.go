package json

import (
	"encoding/json"
	"errors"
	"strings"

	v2net "github.com/v2ray/v2ray-core/common/net"
)

type NetworkList []string

func NewNetworkList(networks []string) NetworkList {
	list := NetworkList(make([]string, len(networks)))
	for idx, network := range networks {
		list[idx] = strings.ToLower(strings.TrimSpace(network))
	}
	return list
}

func (this *NetworkList) UnmarshalJSON(data []byte) error {
	var strList []string
	err := json.Unmarshal(data, &strList)
	if err == nil {
		*this = NewNetworkList(strList)
		return nil
	}

	var str string
	err = json.Unmarshal(data, &str)
	if err == nil {
		strList := strings.Split(str, ",")
		*this = NewNetworkList(strList)
		return nil
	}
	return errors.New("Unknown format of network list: " + string(data))
}

func (this *NetworkList) HasNetwork(network v2net.Network) bool {
	for _, value := range *this {
		if value == string(network) {
			return true
		}
	}
	return false
}
