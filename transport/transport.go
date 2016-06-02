package transport

import "github.com/v2ray/v2ray-core/common/log"

var (
	TCPStreamConfig = TCPConfig{
		ConnectionReuse: false,
	}
)

func ApplyConfig(config *Config) error {
	if config.StreamType == StreamTypeTCP {
		if config.TCPConfig != nil {
			TCPStreamConfig.ConnectionReuse = config.TCPConfig.ConnectionReuse
			if TCPStreamConfig.ConnectionReuse {
				log.Info("Transport: TCP connection reuse enabled.")
			}
		}
	}

	return nil
}
