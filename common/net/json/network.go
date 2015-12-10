package json

import (
	"strings"

	v2net "github.com/v2ray/v2ray-core/common/net"
	serialjson "github.com/v2ray/v2ray-core/common/serial/json"
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
	strlist, err := serialjson.UnmarshalStringList(data)
	if err != nil {
		return err
	}
	*this = NewNetworkList(strlist)
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
