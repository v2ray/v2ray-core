package command_test

import (
	"context"
	"testing"

	"v2ray.com/core"
	"v2ray.com/core/app/dispatcher"
	"v2ray.com/core/app/log"
	. "v2ray.com/core/app/log/command"
	"v2ray.com/core/app/proxyman"
	_ "v2ray.com/core/app/proxyman/inbound"
	_ "v2ray.com/core/app/proxyman/outbound"
	"v2ray.com/core/common/serial"
	. "v2ray.com/ext/assert"
)

func TestLoggerRestart(t *testing.T) {
	assert := With(t)

	v, err := core.New(&core.Config{
		App: []*serial.TypedMessage{
			serial.ToTypedMessage(&log.Config{}),
			serial.ToTypedMessage(&dispatcher.Config{}),
			serial.ToTypedMessage(&proxyman.InboundConfig{}),
			serial.ToTypedMessage(&proxyman.OutboundConfig{}),
		},
	})

	assert(err, IsNil)
	assert(v.Start(), IsNil)

	server := &LoggerServer{
		V: v,
	}
	_, err = server.RestartLogger(context.Background(), &RestartLoggerRequest{})
	assert(err, IsNil)
}
