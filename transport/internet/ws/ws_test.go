package ws_test

import (
	"testing"
	"time"

	"github.com/v2ray/v2ray-core/testing/assert"

	. "github.com/v2ray/v2ray-core/transport/internet/ws"

	v2net "github.com/v2ray/v2ray-core/common/net"
)

func Test_Connect_ws(t *testing.T) {
	assert := assert.On(t)
	(&Config{Pto: "ws", Path: ""}).Apply()
	conn, err := Dial(v2net.AnyIP, v2net.TCPDestination(v2net.DomainAddress("echo.websocket.org"), 80))
	assert.Error(err).IsNil()
	conn.Write([]byte("echo"))
	s := make(chan int)
	go func() {
		buf := make([]byte, 4)
		conn.Read(buf)
		s <- 0
	}()
	<-s
	conn.Close()
}

func Test_Connect_wss(t *testing.T) {
	assert := assert.On(t)
	(&Config{Pto: "wss", Path: ""}).Apply()
	conn, err := Dial(v2net.AnyIP, v2net.TCPDestination(v2net.DomainAddress("echo.websocket.org"), 443))
	assert.Error(err).IsNil()
	conn.Write([]byte("echo"))
	s := make(chan int)
	go func() {
		buf := make([]byte, 4)
		conn.Read(buf)
		s <- 0
	}()
	<-s
	conn.Close()
}

func Test_Connect_ws_guess(t *testing.T) {
	assert := assert.On(t)
	(&Config{Pto: "", Path: ""}).Apply()
	conn, err := Dial(v2net.AnyIP, v2net.TCPDestination(v2net.DomainAddress("echo.websocket.org"), 80))
	assert.Error(err).IsNil()
	conn.Write([]byte("echo"))
	s := make(chan int)
	go func() {
		buf := make([]byte, 4)
		conn.Read(buf)
		s <- 0
	}()
	<-s
	conn.Close()
}

func Test_Connect_wss_guess(t *testing.T) {
	assert := assert.On(t)
	(&Config{Pto: "", Path: ""}).Apply()
	conn, err := Dial(v2net.AnyIP, v2net.TCPDestination(v2net.DomainAddress("echo.websocket.org"), 443))
	assert.Error(err).IsNil()
	conn.Write([]byte("echo"))
	s := make(chan int)
	go func() {
		buf := make([]byte, 4)
		conn.Read(buf)
		s <- 0
	}()
	<-s
	conn.Close()
}

func Test_Connect_wss_guess_fail(t *testing.T) {
	assert := assert.On(t)
	(&Config{Pto: "", Path: ""}).Apply()
	_, err := Dial(v2net.AnyIP, v2net.TCPDestination(v2net.DomainAddress("static.kkdev.org"), 443))
	assert.Error(err).IsNotNil()
}

func Test_Connect_wss_guess_reuse(t *testing.T) {
	assert := assert.On(t)
	(&Config{Pto: "", Path: "", ConnectionReuse: true}).Apply()
	i := 3
	for i != 0 {
		conn, err := Dial(v2net.AnyIP, v2net.TCPDestination(v2net.DomainAddress("echo.websocket.org"), 443))
		assert.Error(err).IsNil()
		conn.Write([]byte("echo"))
		s := make(chan int)
		go func() {
			buf := make([]byte, 4)
			conn.Read(buf)
			s <- 0
		}()
		<-s
		if i == 0 {
			conn.SetDeadline(time.Now())
			conn.SetReadDeadline(time.Now())
			conn.SetWriteDeadline(time.Now())
			conn.SetReusable(false)
		}
		conn.Close()
		i--
	}
}

func Test_listenWSAndDial(t *testing.T) {
	assert := assert.On(t)
	(&Config{Pto: "ws", Path: ""}).Apply()
	listen, err := ListenWS(v2net.DomainAddress("localhost"), 13142)
	assert.Error(err).IsNil()
	go func() {
		conn, err := listen.Accept()
		assert.Error(err).IsNil()
		conn.Close()
	}()
	conn, err := Dial(v2net.AnyIP, v2net.TCPDestination(v2net.DomainAddress("localhost"), 13142))
	assert.Error(err).IsNil()
	conn.Close()
}
