package protocol_test

import (
	"testing"
	"time"

	"v2ray.com/core/common/net"
	. "v2ray.com/core/common/protocol"
	"v2ray.com/core/testing/assert"
)

func TestServerList(t *testing.T) {
	assert := assert.On(t)

	list := NewServerList()
	list.AddServer(NewServerSpec(net.TCPDestination(net.LocalHostIP, net.Port(1)), AlwaysValid()))
	assert.Uint32(list.Size()).Equals(1)
	list.AddServer(NewServerSpec(net.TCPDestination(net.LocalHostIP, net.Port(2)), BeforeTime(time.Now().Add(time.Second))))
	assert.Uint32(list.Size()).Equals(2)

	server := list.GetServer(1)
	assert.Port(server.Destination().Port).Equals(2)
	time.Sleep(2 * time.Second)
	server = list.GetServer(1)
	assert.Pointer(server).IsNil()

	server = list.GetServer(0)
	assert.Port(server.Destination().Port).Equals(1)
}

func TestServerPicker(t *testing.T) {
	assert := assert.On(t)

	list := NewServerList()
	list.AddServer(NewServerSpec(net.TCPDestination(net.LocalHostIP, net.Port(1)), AlwaysValid()))
	list.AddServer(NewServerSpec(net.TCPDestination(net.LocalHostIP, net.Port(2)), BeforeTime(time.Now().Add(time.Second))))
	list.AddServer(NewServerSpec(net.TCPDestination(net.LocalHostIP, net.Port(3)), BeforeTime(time.Now().Add(time.Second))))

	picker := NewRoundRobinServerPicker(list)
	server := picker.PickServer()
	assert.Port(server.Destination().Port).Equals(1)
	server = picker.PickServer()
	assert.Port(server.Destination().Port).Equals(2)
	server = picker.PickServer()
	assert.Port(server.Destination().Port).Equals(3)
	server = picker.PickServer()
	assert.Port(server.Destination().Port).Equals(1)

	time.Sleep(2 * time.Second)
	server = picker.PickServer()
	assert.Port(server.Destination().Port).Equals(1)
	server = picker.PickServer()
	assert.Port(server.Destination().Port).Equals(1)
}
