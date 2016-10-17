package inbound

import ()

func (this *Config) GetDefaultValue() *DefaultConfig {
	if this.GetDefault() == nil {
		return &DefaultConfig{
			AlterId: 32,
			Level:   0,
		}
	}
	return this.Default
}
