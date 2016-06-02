package transport

type StreamType int

type Config struct {
	ConnectionReuse bool
}

func (this *Config) Apply() error {
	if this.ConnectionReuse {
		connectionReuse = true
	}
	return nil
}
