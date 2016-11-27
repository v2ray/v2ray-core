package inbound

func (v *Config) GetDefaultValue() *DefaultConfig {
	if v.GetDefault() == nil {
		return &DefaultConfig{
			AlterId: 32,
			Level:   0,
		}
	}
	return v.Default
}
