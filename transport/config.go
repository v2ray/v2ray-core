package transport

// Config for V2Ray transport layer.
type Config struct {
	ConnectionReuse bool
}

// Apply applies this Config.
func (this *Config) Apply() error {
	if this.ConnectionReuse {
		connectionReuse = true
	}
	return nil
}
