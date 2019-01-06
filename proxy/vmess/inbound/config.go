package inbound

// GetDefaultValue returns default settings of DefaultConfig.
func (c *Config) GetDefaultValue() *DefaultConfig {
	if c.GetDefault() == nil {
		return &DefaultConfig{
			AlterId: 32,
			Level:   0,
		}
	}
	return c.Default
}
