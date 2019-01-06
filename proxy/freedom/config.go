package freedom

func (c *Config) useIP() bool {
	return c.DomainStrategy == Config_USE_IP || c.DomainStrategy == Config_USE_IP4 || c.DomainStrategy == Config_USE_IP6
}
