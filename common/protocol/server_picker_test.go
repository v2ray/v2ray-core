package protocol_test

import (
	"testing"
	"time"

	v2net "github.com/v2ray/v2ray-core/common/net"
	. "github.com/v2ray/v2ray-core/common/protocol"
	"github.com/v2ray/v2ray-core/testing/assert"
)

func TestServerList(t *testing.T) {
	assert := assert.On(t)

	list := NewServerList()
	list.AddServer(NewServerSpec(v2net.TCPDestination(v2net.LocalHostIP, v2net.Port(1)), AlwaysValid()))
	assert.Uint32(list.Size()).Equals(1)
	list.AddServer(NewServerSpec(v2net.TCPDestination(v2net.LocalHostIP, v2net.Port(2)), BeforeTime(time.Now().Add(time.Second))))
	assert.Uint32(list.Size()).Equals(2)

	server := list.GetServer(1)
	assert.Port(server.Destination().Port()).Equals(2)
	time.Sleep(2 * time.Second)
	server = list.GetServer(1)
	assert.Pointer(server).IsNil()

	server = list.GetServer(0)
	assert.Port(server.Destination().Port()).Equals(1)
}

func TestServerPicker(t *testing.T) {
	assert := assert.On(t)

	list := NewServerList()
	list.AddServer(NewServerSpec(v2net.TCPDestination(v2net.LocalHostIP, v2net.Port(1)), AlwaysValid()))
	list.AddServer(NewServerSpec(v2net.TCPDestination(v2net.LocalHostIP, v2net.Port(2)), BeforeTime(time.Now().Add(time.Second))))
	list.AddServer(NewServerSpec(v2net.TCPDestination(v2net.LocalHostIP, v2net.Port(3)), BeforeTime(time.Now().Add(time.Second))))

	picker := NewRoundRobinServerPicker(list)
	server := picker.PickServer()
	assert.Port(server.Destination().Port()).Equals(1)
	server = picker.PickServer()
	assert.Port(server.Destination().Port()).Equals(2)
	server = picker.PickServer()
	assert.Port(server.Destination().Port()).Equals(3)

	time.Sleep(2 * time.Second)
	server = picker.PickServer()
	assert.Port(server.Destination().Port()).Equals(1)
	server = picker.PickServer()
	assert.Port(server.Destination().Port()).Equals(1)
}
