package tcp

type Config struct {
	ConnectionReuse bool
}

func (this *Config) Apply() {
	effectiveConfig = this
}

var (
	effectiveConfig = &Config{
		ConnectionReuse: true,
	}
)
