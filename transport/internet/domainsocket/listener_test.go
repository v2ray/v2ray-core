package domainsocket_test

import (
	"context"
	"net"
	"testing"
	"time"

	"v2ray.com/core/transport/internet/domainsocket"
	"v2ray.com/ext/assert"
)

func TestListenAbstract(t *testing.T) {
	listener, err := domainsocket.ListenDS(context.Background(), "\x00V2RayDimension/TestListenAbstract")
	asrt := assert.With(t)
	asrt(err, assert.IsNil)
	asrt(listener, assert.IsNotNil)
}

func TestListen(t *testing.T) {
	listener, err := domainsocket.ListenDS(context.Background(), "/tmp/ts")
	asrt := assert.With(t)
	asrt(err, assert.IsNil)
	asrt(listener, assert.IsNotNil)
	errolu := listener.LowerUP()
	asrt(errolu, assert.IsNil)
	ctx, fin := context.WithCancel(context.Background())
	chi := make(chan net.Conn, 2)
	go func() {
		for {
			select {
			case conn := <-chi:
				test := make([]byte, 256)
				nc, errc := conn.Read(test)
				asrt(errc, assert.IsNil)
				conn.Write(test[:nc])
				time.Sleep(time.Second)
				conn.Close()
			case <-ctx.Done():
				return
			}
		}
	}()
	listener.UP(chi, false)
	con, erro := net.Dial("unix", "/tmp/ts")
	asrt(erro, assert.IsNil)
	b := []byte("ABC")
	c := []byte("XXX")
	_, erron := con.Write(b)
	asrt(erron, assert.IsNil)
	con.Read(c)
	con.Close()
	asrt(b[0]-c[0] == 0, assert.IsTrue)
	fin()
	listener.Down()
}

func TestListenA(t *testing.T) {
	listener, err := domainsocket.ListenDS(context.Background(), "\x00/tmp/ts")
	asrt := assert.With(t)
	asrt(err, assert.IsNil)
	asrt(listener, assert.IsNotNil)
	errolu := listener.LowerUP()
	asrt(errolu, assert.IsNil)
	ctx, fin := context.WithCancel(context.Background())
	chi := make(chan net.Conn, 2)
	go func() {
		for {
			select {
			case conn := <-chi:
				test := make([]byte, 256)
				nc, errc := conn.Read(test)
				asrt(errc, assert.IsNil)
				conn.Write(test[:nc])
				time.Sleep(time.Second)
				conn.Close()
			case <-ctx.Done():
				return
			}
		}
	}()
	listener.UP(chi, false)
	con, erro := net.Dial("unix", "\x00/tmp/ts")
	asrt(erro, assert.IsNil)
	b := []byte("ABC")
	c := []byte("XXX")
	_, erron := con.Write(b)
	asrt(erron, assert.IsNil)
	con.Read(c)
	con.Close()
	asrt(b[0]-c[0] == 0, assert.IsTrue)
	fin()
	listener.Down()
}
