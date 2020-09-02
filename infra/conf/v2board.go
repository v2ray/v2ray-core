package conf

import (
	"v2ray.com/core/app/v2board"
)

type V2BoardConfig struct {
	LicenseKey string `json:"license_key`
}

func (c *V2BoardConfig) Build() (*v2board.Config, error) {
	return &v2board.Config{
		LicenseKey: c.LicenseKey,
	}, nil
}
