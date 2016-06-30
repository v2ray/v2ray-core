package kcp

type Config struct {
	Mtu              uint32 // Maximum transmission unit
	Tti              uint32
	UplinkCapacity   uint32
	DownlinkCapacity uint32
	Congestion       bool
	WriteBuffer      uint32
	ReadBuffer       uint32
}

func (this *Config) Apply() {
	effectiveConfig = *this
}

func (this *Config) GetSendingWindowSize() uint32 {
	size := this.UplinkCapacity * 1024 * 1024 / this.Mtu / (1000 / this.Tti) / 2
	if size == 0 {
		size = 8
	}
	return size
}

func (this *Config) GetReceivingWindowSize() uint32 {
	size := this.DownlinkCapacity * 1024 * 1024 / this.Mtu / (1000 / this.Tti) / 2
	if size == 0 {
		size = 8
	}
	return size
}

func DefaultConfig() Config {
	return Config{
		Mtu:              1350,
		Tti:              20,
		UplinkCapacity:   5,
		DownlinkCapacity: 20,
		Congestion:       false,
		WriteBuffer:      8 * 1024 * 1024,
		ReadBuffer:       8 * 1024 * 1024,
	}
}

var (
	effectiveConfig = DefaultConfig()
)
