package kcp

type Config struct {
	Mtu        int  // Maximum transmission unit
	Sndwnd     int  // Sending window size
	Rcvwnd     int  // Receiving window size
	Acknodelay bool // Acknoledge without delay
}

func (this *Config) Apply() {
	effectiveConfig = *this
}

func DefaultConfig() Config {
	return Config{
		Mtu:        1350,
		Sndwnd:     1024,
		Rcvwnd:     1024,
		Acknodelay: true,
	}
}

var (
	effectiveConfig = DefaultConfig()
)
