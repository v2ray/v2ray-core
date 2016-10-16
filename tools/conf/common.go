package conf

import (
	"encoding/json"
	"errors"

	"strings"
	v2net "v2ray.com/core/common/net"
)

type Address struct {
	v2net.Address
}

func (this *Address) UnmarshalJSON(data []byte) error {
	var rawStr string
	if err := json.Unmarshal(data, &rawStr); err != nil {
		return err
	}
	this.Address = v2net.ParseAddress(rawStr)

	return nil
}

func (this *Address) Build() *v2net.IPOrDomain {
	if this.Family().IsDomain() {
		return &v2net.IPOrDomain{
			Address: &v2net.IPOrDomain_Domain{
				Domain: this.Domain(),
			},
		}
	}

	return &v2net.IPOrDomain{
		Address: &v2net.IPOrDomain_Ip{
			Ip: []byte(this.IP()),
		},
	}
}

type Network string

func (this Network) Build() v2net.Network {
	return v2net.ParseNetwork(string(this))
}

type NetworkList []Network

func (this *NetworkList) UnmarshalJSON(data []byte) error {
	var strarray []Network
	if err := json.Unmarshal(data, &strarray); err == nil {
		nl := NetworkList(strarray)
		*this = nl
		return nil
	}

	var rawstr Network
	if err := json.Unmarshal(data, &rawstr); err == nil {
		strlist := strings.Split(string(rawstr), ",")
		nl := make([]Network, len(strlist))
		for idx, network := range strlist {
			nl[idx] = Network(network)
		}
		*this = nl
		return nil
	}
	return errors.New("Unknown format of a string list: " + string(data))
}

func (this *NetworkList) Build() *v2net.NetworkList {
	list := new(v2net.NetworkList)
	for _, network := range *this {
		list.Network = append(list.Network, network.Build())
	}
	return list
}
