package udp_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"v2ray.com/core"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/net"
	. "v2ray.com/core/transport/internet/udp"
	"v2ray.com/core/transport/pipe"
	. "v2ray.com/ext/assert"
)

type TestDispatcher struct {
	OnDispatch func(ctx context.Context, dest net.Destination) (*core.Link, error)
}

func (d *TestDispatcher) Dispatch(ctx context.Context, dest net.Destination) (*core.Link, error) {
	return d.OnDispatch(ctx, dest)
}

func (d *TestDispatcher) Start() error {
	return nil
}

func (d *TestDispatcher) Close() error {
	return nil
}

func TestSameDestinationDispatching(t *testing.T) {
	assert := With(t)

	ctx, cancel := context.WithCancel(context.Background())
	uplinkReader, uplinkWriter := pipe.New(pipe.WithSizeLimit(1024))
	downlinkReader, downlinkWriter := pipe.New(pipe.WithSizeLimit(1024))

	go func() {
		for {
			data, err := uplinkReader.ReadMultiBuffer()
			if err != nil {
				break
			}
			err = downlinkWriter.WriteMultiBuffer(data)
			assert(err, IsNil)
		}
	}()

	var count uint32
	td := &TestDispatcher{
		OnDispatch: func(ctx context.Context, dest net.Destination) (*core.Link, error) {
			atomic.AddUint32(&count, 1)
			return &core.Link{Reader: downlinkReader, Writer: uplinkWriter}, nil
		},
	}
	dest := net.UDPDestination(net.LocalHostIP, 53)

	b := buf.New()
	b.WriteBytes('a', 'b', 'c', 'd')

	var msgCount uint32
	dispatcher := NewDispatcher(td, func(ctx context.Context, payload *buf.Buffer) {
		atomic.AddUint32(&msgCount, 1)
	})

	dispatcher.Dispatch(ctx, dest, b)
	for i := 0; i < 5; i++ {
		dispatcher.Dispatch(ctx, dest, b)
	}

	time.Sleep(time.Second)
	cancel()

	assert(count, Equals, uint32(1))
	assert(atomic.LoadUint32(&msgCount), Equals, uint32(6))
}
