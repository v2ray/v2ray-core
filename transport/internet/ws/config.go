package ws

type Config struct {
	ConnectionReuse bool
	Path            string
}

func (this *Config) Apply() {
	effectiveConfig = this
}

var (
	effectiveConfig = &Config{
		ConnectionReuse: true,
		Path:            "",
	}
)
