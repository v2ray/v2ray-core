package conf

import (
	"encoding/json"
	"errors"

	"strings"
	"v2ray.com/core/common/log"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/common/protocol"
)

type StringList []string

func NewStringList(raw []string) *StringList {
	list := StringList(raw)
	return &list
}

func (this StringList) Len() int {
	return len(this)
}

func (this *StringList) UnmarshalJSON(data []byte) error {
	var strarray []string
	if err := json.Unmarshal(data, &strarray); err == nil {
		*this = *NewStringList(strarray)
		return nil
	}

	var rawstr string
	if err := json.Unmarshal(data, &rawstr); err == nil {
		strlist := strings.Split(rawstr, ",")
		*this = *NewStringList(strlist)
		return nil
	}
	return errors.New("Unknown format of a string list: " + string(data))
}

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

func parseIntPort(data []byte) (v2net.Port, error) {
	var intPort uint32
	err := json.Unmarshal(data, &intPort)
	if err != nil {
		return v2net.Port(0), err
	}
	return v2net.PortFromInt(intPort)
}

func parseStringPort(data []byte) (v2net.Port, v2net.Port, error) {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return v2net.Port(0), v2net.Port(0), err
	}
	pair := strings.SplitN(s, "-", 2)
	if len(pair) == 0 {
		return v2net.Port(0), v2net.Port(0), v2net.ErrInvalidPortRange
	}
	if len(pair) == 1 {
		port, err := v2net.PortFromString(pair[0])
		return port, port, err
	}

	fromPort, err := v2net.PortFromString(pair[0])
	if err != nil {
		return v2net.Port(0), v2net.Port(0), err
	}
	toPort, err := v2net.PortFromString(pair[1])
	if err != nil {
		return v2net.Port(0), v2net.Port(0), err
	}
	return fromPort, toPort, nil
}

type PortRange struct {
	From uint32
	To   uint32
}

func (this *PortRange) Build() *v2net.PortRange {
	return &v2net.PortRange{
		From: this.From,
		To:   this.To,
	}
}

// UnmarshalJSON implements encoding/json.Unmarshaler.UnmarshalJSON
func (this *PortRange) UnmarshalJSON(data []byte) error {
	port, err := parseIntPort(data)
	if err == nil {
		this.From = uint32(port)
		this.To = uint32(port)
		return nil
	}

	from, to, err := parseStringPort(data)
	if err == nil {
		this.From = uint32(from)
		this.To = uint32(to)
		if this.From > this.To {
			log.Error("Invalid port range ", this.From, " -> ", this.To)
			return v2net.ErrInvalidPortRange
		}
		return nil
	}

	log.Error("Invalid port range: ", string(data))
	return v2net.ErrInvalidPortRange
}

type User struct {
	EmailString string `json:"email"`
	LevelByte   byte   `json:"level"`
}

func (this *User) Build() *protocol.User {
	return &protocol.User{
		Email: this.EmailString,
		Level: uint32(this.LevelByte),
	}
}
