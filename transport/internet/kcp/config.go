package kcp

type Config struct {
	Mtu              int // Maximum transmission unit
	Tti              int
	UplinkCapacity   int
	DownlinkCapacity int
	Congestion       bool
}

func (this *Config) Apply() {
	effectiveConfig = *this
}

func (this *Config) GetSendingWindowSize() int {
	return this.UplinkCapacity * 1024 * 1024 / this.Mtu / (1000 / this.Tti)
}

func (this *Config) GetReceivingWindowSize() int {
	return this.DownlinkCapacity * 1024 * 1024 / this.Mtu / (1000 / this.Tti)
}

func DefaultConfig() Config {
	return Config{
		Mtu:              1350,
		Tti:              20,
		UplinkCapacity:   5,
		DownlinkCapacity: 20,
		Congestion:       false,
	}
}

var (
	effectiveConfig = DefaultConfig()
)
