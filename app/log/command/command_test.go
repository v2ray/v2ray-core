package command_test

import (
	"context"
	"testing"

	"v2ray.com/core"
	"v2ray.com/core/app/log"
	. "v2ray.com/core/app/log/command"
	"v2ray.com/core/common/serial"
	. "v2ray.com/ext/assert"
)

func TestLoggerRestart(t *testing.T) {
	assert := With(t)

	v, err := core.New(&core.Config{
		App: []*serial.TypedMessage{
			serial.ToTypedMessage(&log.Config{}),
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
