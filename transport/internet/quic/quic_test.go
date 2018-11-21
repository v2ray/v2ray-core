package quic_test

import (
	"context"
	"crypto/rand"
	"testing"
	"time"

	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/common/protocol/tls/cert"
	"v2ray.com/core/common/serial"
	"v2ray.com/core/testing/servers/udp"
	"v2ray.com/core/transport/internet"
	"v2ray.com/core/transport/internet/headers/wireguard"
	"v2ray.com/core/transport/internet/quic"
	"v2ray.com/core/transport/internet/tls"
	. "v2ray.com/ext/assert"
)

func TestQuicConnection(t *testing.T) {
	assert := With(t)

	port := udp.PickPort()

	listener, err := quic.Listen(context.Background(), net.LocalHostIP, port, &internet.MemoryStreamConfig{
		ProtocolName:     "quic",
		ProtocolSettings: &quic.Config{},
		SecurityType:     "tls",
		SecuritySettings: &tls.Config{
			Certificate: []*tls.Certificate{tls.ParseCertificate(cert.MustGenerate(nil, cert.DNSNames("www.v2ray.com"), cert.CommonName("www.v2ray.com")))},
		},
	}, func(conn internet.Connection) {
		go func() {
			defer conn.Close()

			b := buf.New()
			defer b.Release()

			for {
				b.Clear()
				if _, err := b.ReadFrom(conn); err != nil {
					return
				}
				nBytes, err := conn.Write(b.Bytes())
				assert(err, IsNil)
				assert(int32(nBytes), Equals, b.Len())
			}
		}()
	})
	assert(err, IsNil)

	defer listener.Close()

	time.Sleep(time.Second)

	dctx := context.Background()
	conn, err := quic.Dial(dctx, net.TCPDestination(net.LocalHostIP, port), &internet.MemoryStreamConfig{
		ProtocolName:     "quic",
		ProtocolSettings: &quic.Config{},
		SecurityType:     "tls",
		SecuritySettings: &tls.Config{
			ServerName:    "www.v2ray.com",
			AllowInsecure: true,
		},
	})
	assert(err, IsNil)
	defer conn.Close()

	const N = 1024
	b1 := make([]byte, N)
	common.Must2(rand.Read(b1))
	b2 := buf.New()

	nBytes, err := conn.Write(b1)
	assert(nBytes, Equals, N)
	assert(err, IsNil)

	b2.Clear()
	common.Must2(b2.ReadFullFrom(conn, N))
	assert(b2.Bytes(), Equals, b1)

	nBytes, err = conn.Write(b1)
	assert(nBytes, Equals, N)
	assert(err, IsNil)

	b2.Clear()
	common.Must2(b2.ReadFullFrom(conn, N))
	assert(b2.Bytes(), Equals, b1)
}

func TestQuicConnectionWithoutTLS(t *testing.T) {
	assert := With(t)

	port := udp.PickPort()

	listener, err := quic.Listen(context.Background(), net.LocalHostIP, port, &internet.MemoryStreamConfig{
		ProtocolName:     "quic",
		ProtocolSettings: &quic.Config{},
	}, func(conn internet.Connection) {
		go func() {
			defer conn.Close()

			b := buf.New()
			defer b.Release()

			for {
				b.Clear()
				if _, err := b.ReadFrom(conn); err != nil {
					return
				}
				nBytes, err := conn.Write(b.Bytes())
				assert(err, IsNil)
				assert(int32(nBytes), Equals, b.Len())
			}
		}()
	})
	assert(err, IsNil)

	defer listener.Close()

	time.Sleep(time.Second)

	dctx := context.Background()
	conn, err := quic.Dial(dctx, net.TCPDestination(net.LocalHostIP, port), &internet.MemoryStreamConfig{
		ProtocolName:     "quic",
		ProtocolSettings: &quic.Config{},
	})
	assert(err, IsNil)
	defer conn.Close()

	const N = 1024
	b1 := make([]byte, N)
	common.Must2(rand.Read(b1))
	b2 := buf.New()

	nBytes, err := conn.Write(b1)
	assert(nBytes, Equals, N)
	assert(err, IsNil)

	b2.Clear()
	common.Must2(b2.ReadFullFrom(conn, N))
	assert(b2.Bytes(), Equals, b1)

	nBytes, err = conn.Write(b1)
	assert(nBytes, Equals, N)
	assert(err, IsNil)

	b2.Clear()
	common.Must2(b2.ReadFullFrom(conn, N))
	assert(b2.Bytes(), Equals, b1)
}

func TestQuicConnectionAuthHeader(t *testing.T) {
	assert := With(t)

	port := udp.PickPort()

	listener, err := quic.Listen(context.Background(), net.LocalHostIP, port, &internet.MemoryStreamConfig{
		ProtocolName: "quic",
		ProtocolSettings: &quic.Config{
			Header: serial.ToTypedMessage(&wireguard.WireguardConfig{}),
			Key:    "abcd",
			Security: &protocol.SecurityConfig{
				Type: protocol.SecurityType_AES128_GCM,
			},
		},
	}, func(conn internet.Connection) {
		go func() {
			defer conn.Close()

			b := buf.New()
			defer b.Release()

			for {
				b.Clear()
				if _, err := b.ReadFrom(conn); err != nil {
					return
				}
				nBytes, err := conn.Write(b.Bytes())
				assert(err, IsNil)
				assert(int32(nBytes), Equals, b.Len())
			}
		}()
	})
	assert(err, IsNil)

	defer listener.Close()

	time.Sleep(time.Second)

	dctx := context.Background()
	conn, err := quic.Dial(dctx, net.TCPDestination(net.LocalHostIP, port), &internet.MemoryStreamConfig{
		ProtocolName: "quic",
		ProtocolSettings: &quic.Config{
			Header: serial.ToTypedMessage(&wireguard.WireguardConfig{}),
			Key:    "abcd",
			Security: &protocol.SecurityConfig{
				Type: protocol.SecurityType_AES128_GCM,
			},
		},
	})
	assert(err, IsNil)
	defer conn.Close()

	const N = 1024
	b1 := make([]byte, N)
	common.Must2(rand.Read(b1))
	b2 := buf.New()

	nBytes, err := conn.Write(b1)
	assert(nBytes, Equals, N)
	assert(err, IsNil)

	b2.Clear()
	common.Must2(b2.ReadFullFrom(conn, N))
	assert(b2.Bytes(), Equals, b1)

	nBytes, err = conn.Write(b1)
	assert(nBytes, Equals, N)
	assert(err, IsNil)

	b2.Clear()
	common.Must2(b2.ReadFullFrom(conn, N))
	assert(b2.Bytes(), Equals, b1)
}
