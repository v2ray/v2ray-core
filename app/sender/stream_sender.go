package sender

import (
	"v2ray.com/core/app"
	"v2ray.com/core/common/net"
	"v2ray.com/core/transport/internet"
)

type StreamSender struct {
	config *StreamSenderConfig
}

func NewStreamSender(config *StreamSenderConfig, space app.Space) *StreamSender {
	return &StreamSender{
		config: config,
	}
}

func (s *StreamSender) SendTo(destination net.Destination) (internet.Connection, error) {
	src := s.config.Via.AsAddress()
	return internet.Dial(src, destination, internet.DialerOptions{
		Stream: s.config.StreamSettings,
		Proxy:  s.config.ProxySettings,
	})
}
