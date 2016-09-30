package ws_test

import (
	"crypto/tls"
	"testing"
	"time"

	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/testing/assert"
	"v2ray.com/core/transport/internet"
	. "v2ray.com/core/transport/internet/ws"
)

func Test_Connect_ws(t *testing.T) {
	assert := assert.On(t)
	(&Config{Path: ""}).Apply()
	conn, err := Dial(v2net.AnyIP, v2net.TCPDestination(v2net.DomainAddress("echo.websocket.org"), 80), internet.DialerOptions{})
	assert.Error(err).IsNil()
	conn.Write([]byte("echo"))
	s := make(chan int)
	go func() {
		buf := make([]byte, 4)
		conn.Read(buf)
		str := string(buf)
		if str != "echo" {
			assert.Fail("Data mismatch")
		}
		s <- 0
	}()
	<-s
	conn.Close()
}

func Test_Connect_wss(t *testing.T) {
	assert := assert.On(t)
	(&Config{Path: ""}).Apply()
	conn, err := Dial(v2net.AnyIP, v2net.TCPDestination(v2net.DomainAddress("echo.websocket.org"), 443), internet.DialerOptions{
		Stream: &internet.StreamSettings{
			Security: internet.StreamSecurityTypeTLS,
		},
	})
	assert.Error(err).IsNil()
	conn.Write([]byte("echo"))
	s := make(chan int)
	go func() {
		buf := make([]byte, 4)
		conn.Read(buf)
		str := string(buf)
		if str != "echo" {
			assert.Fail("Data mismatch")
		}
		s <- 0
	}()
	<-s
	conn.Close()
}

func Test_Connect_wss_1_nil(t *testing.T) {
	assert := assert.On(t)
	(&Config{Path: ""}).Apply()
	conn, err := Dial(nil, v2net.TCPDestination(v2net.DomainAddress("echo.websocket.org"), 443), internet.DialerOptions{
		Stream: &internet.StreamSettings{
			Security: internet.StreamSecurityTypeTLS,
		},
	})
	assert.Error(err).IsNil()
	conn.Write([]byte("echo"))
	s := make(chan int)
	go func() {
		buf := make([]byte, 4)
		conn.Read(buf)
		str := string(buf)
		if str != "echo" {
			assert.Fail("Data mismatch")
		}
		s <- 0
	}()
	<-s
	conn.Close()
}

func Test_Connect_ws_guess(t *testing.T) {
	assert := assert.On(t)
	(&Config{Path: ""}).Apply()
	conn, err := Dial(v2net.AnyIP, v2net.TCPDestination(v2net.DomainAddress("echo.websocket.org"), 80), internet.DialerOptions{})
	assert.Error(err).IsNil()
	conn.Write([]byte("echo"))
	s := make(chan int)
	go func() {
		buf := make([]byte, 4)
		conn.Read(buf)
		str := string(buf)
		if str != "echo" {
			assert.Fail("Data mismatch")
		}
		s <- 0
	}()
	<-s
	conn.Close()
}

func Test_Connect_wss_guess(t *testing.T) {
	assert := assert.On(t)
	(&Config{Path: ""}).Apply()
	conn, err := Dial(v2net.AnyIP, v2net.TCPDestination(v2net.DomainAddress("echo.websocket.org"), 443), internet.DialerOptions{
		Stream: &internet.StreamSettings{
			Security: internet.StreamSecurityTypeTLS,
		},
	})
	assert.Error(err).IsNil()
	conn.Write([]byte("echo"))
	s := make(chan int)
	go func() {
		buf := make([]byte, 4)
		conn.Read(buf)
		str := string(buf)
		if str != "echo" {
			assert.Fail("Data mismatch")
		}
		s <- 0
	}()
	<-s
	conn.Close()
}

func Test_Connect_wss_guess_fail(t *testing.T) {
	assert := assert.On(t)
	(&Config{Path: ""}).Apply()
	_, err := Dial(v2net.AnyIP, v2net.TCPDestination(v2net.DomainAddress("static.kkdev.org"), 443), internet.DialerOptions{
		Stream: &internet.StreamSettings{
			Security: internet.StreamSecurityTypeTLS,
		},
	})
	assert.Error(err).IsNotNil()
}

func Test_Connect_wss_guess_reuse(t *testing.T) {
	assert := assert.On(t)
	(&Config{Path: "", ConnectionReuse: true}).Apply()
	i := 3
	for i != 0 {
		conn, err := Dial(v2net.AnyIP, v2net.TCPDestination(v2net.DomainAddress("echo.websocket.org"), 443), internet.DialerOptions{
			Stream: &internet.StreamSettings{
				Security: internet.StreamSecurityTypeTLS,
			},
		})
		assert.Error(err).IsNil()
		conn.Write([]byte("echo"))
		s := make(chan int)
		go func() {
			buf := make([]byte, 4)
			conn.Read(buf)
			str := string(buf)
			if str != "echo" {
				assert.Fail("Data mismatch")
			}
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
	(&Config{Path: "ws"}).Apply()
	listen, err := ListenWS(v2net.DomainAddress("localhost"), 13142, internet.ListenOptions{})
	assert.Error(err).IsNil()
	go func() {
		conn, err := listen.Accept()
		assert.Error(err).IsNil()
		conn.Close()
		conn, err = listen.Accept()
		assert.Error(err).IsNil()
		conn.Close()
		conn, err = listen.Accept()
		assert.Error(err).IsNil()
		conn.Close()
		listen.Close()
	}()
	conn, err := Dial(v2net.AnyIP, v2net.TCPDestination(v2net.DomainAddress("localhost"), 13142), internet.DialerOptions{})
	assert.Error(err).IsNil()
	conn.Close()
	<-time.After(time.Second * 5)
	conn, err = Dial(v2net.AnyIP, v2net.TCPDestination(v2net.DomainAddress("localhost"), 13142), internet.DialerOptions{})
	assert.Error(err).IsNil()
	conn.Close()
	<-time.After(time.Second * 15)
	conn, err = Dial(v2net.AnyIP, v2net.TCPDestination(v2net.DomainAddress("localhost"), 13142), internet.DialerOptions{})
	assert.Error(err).IsNil()
	conn.Close()
}

func Test_listenWSAndDial_TLS(t *testing.T) {
	assert := assert.On(t)
	go func() {
		<-time.After(time.Second * 5)
		assert.Fail("Too slow")
	}()
	(&Config{Path: "wss", ConnectionReuse: true}).Apply()

	listen, err := ListenWS(v2net.DomainAddress("localhost"), 13143, internet.ListenOptions{
		Stream: &internet.StreamSettings{
			Security: internet.StreamSecurityTypeTLS,
			TLSSettings: &internet.TLSSettings{
				AllowInsecure: true,
				Certs:         LoadTestCert(assert),
			},
		},
	})
	assert.Error(err).IsNil()
	go func() {
		conn, err := listen.Accept()
		assert.Error(err).IsNil()
		conn.Close()
		listen.Close()
	}()
	conn, err := Dial(v2net.AnyIP, v2net.TCPDestination(v2net.DomainAddress("localhost"), 13143), internet.DialerOptions{
		Stream: &internet.StreamSettings{
			Security: internet.StreamSecurityTypeTLS,
			TLSSettings: &internet.TLSSettings{
				AllowInsecure: true,
				Certs:         LoadTestCert(assert),
			},
		},
	})
	assert.Error(err).IsNil()
	conn.Close()
}

func LoadTestCert(assert *assert.Assert) []tls.Certificate {
	cert, err := tls.LoadX509KeyPair("./../../../testing/tls/cert.pem", "./../../../testing/tls/key.pem")
	assert.Error(err).IsNil()
	return []tls.Certificate{cert}
}
