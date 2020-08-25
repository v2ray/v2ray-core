package protocol_test

import (
	"testing"
	"time"

	"v2ray.com/core/common/net"
	. "v2ray.com/core/common/protocol"
)

func TestServerList(t *testing.T) {
	list := NewServerList()
	list.AddServer(NewServerSpec(net.TCPDestination(net.LocalHostIP, net.Port(1)), AlwaysValid()))
	if list.Size() != 1 {
		t.Error("list size: ", list.Size())
	}
	list.AddServer(NewServerSpec(net.TCPDestination(net.LocalHostIP, net.Port(2)), BeforeTime(time.Now().Add(time.Second))))
	if list.Size() != 2 {
		t.Error("list.size: ", list.Size())
	}

	server := list.GetServer(1)
	if server.Destination().Port != 2 {
		t.Error("server: ", server.Destination())
	}
	time.Sleep(2 * time.Second)
	server = list.GetServer(1)
	if server != nil {
		t.Error("server: ", server)
	}

	server = list.GetServer(0)
	if server.Destination().Port != 1 {
		t.Error("server: ", server.Destination())
	}
}

func TestServerPicker(t *testing.T) {
	list := NewServerList()
	list.AddServer(NewServerSpec(net.TCPDestination(net.LocalHostIP, net.Port(1)), AlwaysValid()))
	list.AddServer(NewServerSpec(net.TCPDestination(net.LocalHostIP, net.Port(2)), BeforeTime(time.Now().Add(time.Second))))
	list.AddServer(NewServerSpec(net.TCPDestination(net.LocalHostIP, net.Port(3)), BeforeTime(time.Now().Add(time.Second))))

	picker := NewRoundRobinServerPicker(list)
	server := picker.PickServer()
	if server.Destination().Port != 1 {
		t.Error("server: ", server.Destination())
	}
	server = picker.PickServer()
	if server.Destination().Port != 2 {
		t.Error("server: ", server.Destination())
	}
	server = picker.PickServer()
	if server.Destination().Port != 3 {
		t.Error("server: ", server.Destination())
	}
	server = picker.PickServer()
	if server.Destination().Port != 1 {
		t.Error("server: ", server.Destination())
	}

	time.Sleep(2 * time.Second)
	server = picker.PickServer()
	if server.Destination().Port != 1 {
		t.Error("server: ", server.Destination())
	}
	server = picker.PickServer()
	if server.Destination().Port != 1 {
		t.Error("server: ", server.Destination())
	}
}
