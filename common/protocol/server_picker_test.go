package protocol_test

import (
	"testing"
	"time"

	"v2ray.com/core/common/net"
	. "v2ray.com/core/common/protocol"
	. "v2ray.com/ext/assert"
)

func TestServerList(t *testing.T) {
	assert := With(t)

	list := NewServerList()
	list.AddServer(NewServerSpec(net.TCPDestination(net.LocalHostIP, net.Port(1)), AlwaysValid()))
	assert(list.Size(), Equals, uint32(1))
	list.AddServer(NewServerSpec(net.TCPDestination(net.LocalHostIP, net.Port(2)), BeforeTime(time.Now().Add(time.Second))))
	assert(list.Size(), Equals, uint32(2))

	server := list.GetServer(1)
	assert(server.Destination().Port, Equals, net.Port(2))
	time.Sleep(2 * time.Second)
	server = list.GetServer(1)
	assert(server, IsNil)

	server = list.GetServer(0)
	assert(server.Destination().Port, Equals, net.Port(1))
}

func TestServerPicker(t *testing.T) {
	assert := With(t)

	list := NewServerList()
	list.AddServer(NewServerSpec(net.TCPDestination(net.LocalHostIP, net.Port(1)), AlwaysValid()))
	list.AddServer(NewServerSpec(net.TCPDestination(net.LocalHostIP, net.Port(2)), BeforeTime(time.Now().Add(time.Second))))
	list.AddServer(NewServerSpec(net.TCPDestination(net.LocalHostIP, net.Port(3)), BeforeTime(time.Now().Add(time.Second))))

	picker := NewRoundRobinServerPicker(list)
	server := picker.PickServer()
	assert(server.Destination().Port, Equals, net.Port(1))
	server = picker.PickServer()
	assert(server.Destination().Port, Equals, net.Port(2))
	server = picker.PickServer()
	assert(server.Destination().Port, Equals, net.Port(3))
	server = picker.PickServer()
	assert(server.Destination().Port, Equals, net.Port(1))

	time.Sleep(2 * time.Second)
	server = picker.PickServer()
	assert(server.Destination().Port, Equals, net.Port(1))
	server = picker.PickServer()
	assert(server.Destination().Port, Equals, net.Port(1))
}
