package transport

type StreamType int

const (
	StreamTypeTCP = StreamType(0)
)

type TCPConfig struct {
	ConnectionReuse bool
}

type Config struct {
	StreamType StreamType
	TCPConfig  *TCPConfig
}
