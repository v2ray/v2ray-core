package kcp

import (
	"context"
	"crypto/tls"
	"io"
	"sync/atomic"

	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/dice"
	"v2ray.com/core/common/net"
	"v2ray.com/core/transport/internet"
	v2tls "v2ray.com/core/transport/internet/tls"
)

var (
	globalConv = uint32(dice.RollUint16())
)

func fetchInput(ctx context.Context, input io.Reader, reader PacketReader, conn *Connection) {
	payload := buf.New()
	defer payload.Release()

	for {
		err := payload.Reset(buf.ReadFrom(input))
		if err != nil {
			payload.Release()
			return
		}
		segments := reader.Read(payload.Bytes())
		if len(segments) > 0 {
			conn.Input(segments)
		}
	}
}

func DialKCP(ctx context.Context, dest net.Destination) (internet.Connection, error) {
	dest.Network = net.Network_UDP
	newError("dialing mKCP to ", dest).WriteToLog()

	src := internet.DialerSourceFromContext(ctx)
	rawConn, err := internet.DialSystem(ctx, src, dest)
	if err != nil {
		return nil, newError("failed to dial to dest: ", err).AtWarning().Base(err)
	}

	kcpSettings := internet.TransportSettingsFromContext(ctx).(*Config)

	header, err := kcpSettings.GetPackerHeader()
	if err != nil {
		return nil, newError("failed to create packet header").Base(err)
	}
	security, err := kcpSettings.GetSecurity()
	if err != nil {
		return nil, newError("failed to create security").Base(err)
	}
	reader := &KCPPacketReader{
		Header:   header,
		Security: security,
	}
	writer := &KCPPacketWriter{
		Header:   header,
		Security: security,
		Writer:   rawConn,
	}

	conv := uint16(atomic.AddUint32(&globalConv, 1))
	session := NewConnection(ConnMetadata{
		LocalAddr:    rawConn.LocalAddr(),
		RemoteAddr:   rawConn.RemoteAddr(),
		Conversation: conv,
	}, writer, rawConn, kcpSettings)

	go fetchInput(ctx, rawConn, reader, session)

	var iConn internet.Connection = session

	if config := v2tls.ConfigFromContext(ctx, v2tls.WithDestination(dest)); config != nil {
		tlsConn := tls.Client(iConn, config.GetTLSConfig())
		iConn = tlsConn
	}

	return iConn, nil
}

func init() {
	common.Must(internet.RegisterTransportDialer(internet.TransportProtocol_MKCP, DialKCP))
}
