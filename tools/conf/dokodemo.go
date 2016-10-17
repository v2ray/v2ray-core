package conf

import (
	"v2ray.com/core/common/loader"
	"v2ray.com/core/proxy/dokodemo"
)

type DokodemoConfig struct {
	Host         *Address     `json:"address"`
	PortValue    uint16       `json:"port"`
	NetworkList  *NetworkList `json:"network"`
	TimeoutValue uint32       `json:"timeout"`
	Redirect     bool         `json:"followRedirect"`
}

func (this *DokodemoConfig) Build() (*loader.TypedSettings, error) {
	config := new(dokodemo.Config)
	if this.Host != nil {
		config.Address = this.Host.Build()
	}
	config.Port = uint32(this.PortValue)
	config.NetworkList = this.NetworkList.Build()
	config.Timeout = this.TimeoutValue
	config.FollowRedirect = this.Redirect
	return loader.NewTypedSettings(config), nil
}
