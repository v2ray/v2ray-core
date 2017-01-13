package dokodemo_test

import (
	"net"
	"testing"

	"context"

	"v2ray.com/core/app"
	"v2ray.com/core/app/dispatcher"
	_ "v2ray.com/core/app/dispatcher/impl"
	"v2ray.com/core/app/proxyman"
	_ "v2ray.com/core/app/proxyman/outbound"
	"v2ray.com/core/common/dice"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/proxy"
	. "v2ray.com/core/proxy/dokodemo"
	"v2ray.com/core/proxy/freedom"
	"v2ray.com/core/testing/assert"
	"v2ray.com/core/testing/servers/tcp"
	"v2ray.com/core/testing/servers/udp"
	"v2ray.com/core/transport/internet"
	_ "v2ray.com/core/transport/internet/tcp"
)

func TestDokodemoTCP(t *testing.T) {
	assert := assert.On(t)

	tcpServer := &tcp.Server{
		MsgProcessor: func(data []byte) []byte {
			buffer := make([]byte, 0, 2048)
			buffer = append(buffer, []byte("Processed: ")...)
			buffer = append(buffer, data...)
			return buffer
		},
	}
	_, err := tcpServer.Start()
	assert.Error(err).IsNil()

	defer tcpServer.Close()

	space := app.NewSpace()
	ctx := app.ContextWithSpace(context.Background(), space)
	app.AddApplicationToSpace(ctx, new(dispatcher.Config))
	app.AddApplicationToSpace(ctx, new(proxyman.OutboundConfig))

	ohm := proxyman.OutboundHandlerManagerFromSpace(space)
	freedom, err := freedom.New(proxy.ContextWithOutboundMeta(ctx, &proxy.OutboundHandlerMeta{
		Address: v2net.LocalHostIP,
		StreamSettings: &internet.StreamConfig{
			Protocol: internet.TransportProtocol_TCP,
		},
	}), &freedom.Config{})
	assert.Error(err).IsNil()
	ohm.SetDefaultHandler(freedom)

	data2Send := "Data to be sent to remote."

	port := v2net.Port(dice.Roll(20000) + 10000)

	ctx = proxy.ContextWithInboundMeta(ctx, &proxy.InboundHandlerMeta{
		Address: v2net.LocalHostIP,
		Port:    port,
		StreamSettings: &internet.StreamConfig{
			Protocol: internet.TransportProtocol_TCP,
		}})

	dokodemo, err := NewDokodemoDoor(ctx, &Config{
		Address:     v2net.NewIPOrDomain(v2net.LocalHostIP),
		Port:        uint32(tcpServer.Port),
		NetworkList: v2net.Network_TCP.AsList(),
		Timeout:     600,
	})
	assert.Error(err).IsNil()
	defer dokodemo.Close()

	assert.Error(space.Initialize()).IsNil()

	err = dokodemo.Start()
	assert.Error(err).IsNil()
	assert.Port(port).Equals(dokodemo.Port())

	tcpClient, err := net.DialTCP("tcp", nil, &net.TCPAddr{
		IP:   []byte{127, 0, 0, 1},
		Port: int(port),
		Zone: "",
	})
	assert.Error(err).IsNil()

	tcpClient.Write([]byte(data2Send))
	tcpClient.CloseWrite()

	response := make([]byte, 1024)
	nBytes, err := tcpClient.Read(response)
	assert.Error(err).IsNil()
	tcpClient.Close()

	assert.String("Processed: " + data2Send).Equals(string(response[:nBytes]))
}

func TestDokodemoUDP(t *testing.T) {
	assert := assert.On(t)

	udpServer := &udp.Server{
		MsgProcessor: func(data []byte) []byte {
			buffer := make([]byte, 0, 2048)
			buffer = append(buffer, []byte("Processed: ")...)
			buffer = append(buffer, data...)
			return buffer
		},
	}
	_, err := udpServer.Start()
	assert.Error(err).IsNil()

	defer udpServer.Close()

	space := app.NewSpace()
	ctx := app.ContextWithSpace(context.Background(), space)
	app.AddApplicationToSpace(ctx, new(dispatcher.Config))
	app.AddApplicationToSpace(ctx, new(proxyman.OutboundConfig))

	ohm := proxyman.OutboundHandlerManagerFromSpace(space)
	freedom, err := freedom.New(proxy.ContextWithOutboundMeta(ctx, &proxy.OutboundHandlerMeta{
		Address: v2net.AnyIP,
		StreamSettings: &internet.StreamConfig{
			Protocol: internet.TransportProtocol_TCP,
		},
	}), &freedom.Config{})
	assert.Error(err).IsNil()
	ohm.SetDefaultHandler(freedom)

	data2Send := "Data to be sent to remote."

	port := v2net.Port(dice.Roll(20000) + 10000)

	ctx = proxy.ContextWithInboundMeta(ctx, &proxy.InboundHandlerMeta{
		Address: v2net.LocalHostIP,
		Port:    port,
		StreamSettings: &internet.StreamConfig{
			Protocol: internet.TransportProtocol_TCP,
		}})

	dokodemo, err := NewDokodemoDoor(ctx, &Config{
		Address:     v2net.NewIPOrDomain(v2net.LocalHostIP),
		Port:        uint32(udpServer.Port),
		NetworkList: v2net.Network_UDP.AsList(),
		Timeout:     600,
	})
	assert.Error(err).IsNil()
	defer dokodemo.Close()

	assert.Error(space.Initialize()).IsNil()

	err = dokodemo.Start()
	assert.Error(err).IsNil()
	assert.Port(port).Equals(dokodemo.Port())

	udpClient, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   []byte{127, 0, 0, 1},
		Port: int(port),
		Zone: "",
	})
	assert.Error(err).IsNil()
	defer udpClient.Close()

	udpClient.Write([]byte(data2Send))

	response := make([]byte, 1024)
	nBytes, addr, err := udpClient.ReadFromUDP(response)
	assert.Error(err).IsNil()
	assert.IP(addr.IP).Equals(v2net.LocalHostIP.IP())
	assert.Bytes(response[:nBytes]).Equals([]byte("Processed: " + data2Send))
}
