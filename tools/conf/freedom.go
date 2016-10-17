package conf

import (
	"strings"

	"v2ray.com/core/common/loader"
	"v2ray.com/core/proxy/freedom"
)

type FreedomConfig struct {
	DomainStrategy string `json:"domainStrategy"`
	Timeout        uint32 `json:"timeout"`
}

func (this *FreedomConfig) Build() (*loader.TypedSettings, error) {
	config := new(freedom.Config)
	config.DomainStrategy = freedom.Config_AS_IS
	domainStrategy := strings.ToLower(this.DomainStrategy)
	if domainStrategy == "useip" || domainStrategy == "use_ip" {
		config.DomainStrategy = freedom.Config_USE_IP
	}
	config.Timeout = this.Timeout
	return loader.NewTypedSettings(config), nil
}
