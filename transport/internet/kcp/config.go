package kcp

type Config struct {
	Mtu              uint32 // Maximum transmission unit
	Tti              uint32
	UplinkCapacity   uint32
	DownlinkCapacity uint32
	Congestion       bool
	WriteBuffer      int
}

func (this *Config) Apply() {
	effectiveConfig = *this
}

func (this *Config) GetSendingWindowSize() uint32 {
	return this.UplinkCapacity * 1024 * 1024 / this.Mtu / (1000 / this.Tti)
}

func (this *Config) GetReceivingWindowSize() uint32 {
	return this.DownlinkCapacity * 1024 * 1024 / this.Mtu / (1000 / this.Tti)
}

func DefaultConfig() Config {
	return Config{
		Mtu:              1350,
		Tti:              20,
		UplinkCapacity:   5,
		DownlinkCapacity: 20,
		Congestion:       false,
		WriteBuffer:      8 * 1024 * 1024,
	}
}

var (
	effectiveConfig = DefaultConfig()
)
