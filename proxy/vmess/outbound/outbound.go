package outbound

import (
	"context"
	"runtime"
	"time"

	"v2ray.com/core/app"
	"v2ray.com/core/app/log"
	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/errors"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/common/retry"
	"v2ray.com/core/common/signal"
	"v2ray.com/core/proxy"
	"v2ray.com/core/proxy/vmess"
	"v2ray.com/core/proxy/vmess/encoding"
	"v2ray.com/core/transport/internet"
	"v2ray.com/core/transport/ray"
)

// Handler is an outbound connection handler for VMess protocol.
type Handler struct {
	serverList   *protocol.ServerList
	serverPicker protocol.ServerPicker
}

func New(ctx context.Context, config *Config) (*Handler, error) {
	space := app.SpaceFromContext(ctx)
	if space == nil {
		return nil, errors.New("VMess|Outbound: No space in context.")
	}

	serverList := protocol.NewServerList()
	for _, rec := range config.Receiver {
		serverList.AddServer(protocol.NewServerSpecFromPB(*rec))
	}
	handler := &Handler{
		serverList:   serverList,
		serverPicker: protocol.NewRoundRobinServerPicker(serverList),
	}

	return handler, nil
}

// Process implements proxy.Outbound.Process().
func (v *Handler) Process(ctx context.Context, outboundRay ray.OutboundRay, dialer proxy.Dialer) error {
	var rec *protocol.ServerSpec
	var conn internet.Connection

	err := retry.ExponentialBackoff(5, 200).On(func() error {
		rec = v.serverPicker.PickServer()
		rawConn, err := dialer.Dial(ctx, rec.Destination())
		if err != nil {
			return err
		}
		conn = rawConn

		return nil
	})
	if err != nil {
		return errors.New("VMess|Outbound: Failed to find an available destination.").Base(err).AtWarning()
	}
	defer conn.Close()

	target, ok := proxy.TargetFromContext(ctx)
	if !ok {
		return errors.New("VMess|Outbound: Target not specified.")
	}
	log.Info("VMess|Outbound: Tunneling request to ", target, " via ", rec.Destination())

	command := protocol.RequestCommandTCP
	if target.Network == net.Network_UDP {
		command = protocol.RequestCommandUDP
	}
	request := &protocol.RequestHeader{
		Version: encoding.Version,
		User:    rec.PickUser(),
		Command: command,
		Address: target.Address,
		Port:    target.Port,
		Option:  protocol.RequestOptionChunkStream,
	}

	rawAccount, err := request.User.GetTypedAccount()
	if err != nil {
		return errors.New("VMess|Outbound: Failed to get user account.").Base(err).AtWarning()
	}
	account := rawAccount.(*vmess.InternalAccount)
	request.Security = account.Security

	if request.Security.Is(protocol.SecurityType_AES128_GCM) || request.Security.Is(protocol.SecurityType_NONE) || request.Security.Is(protocol.SecurityType_CHACHA20_POLY1305) {
		request.Option.Set(protocol.RequestOptionChunkMasking)
	}

	conn.SetReusable(true)
	if conn.Reusable() { // Conn reuse may be disabled on transportation layer
		request.Option.Set(protocol.RequestOptionConnectionReuse)
	}

	input := outboundRay.OutboundInput()
	output := outboundRay.OutboundOutput()

	session := encoding.NewClientSession(protocol.DefaultIDHash)

	ctx, timer := signal.CancelAfterInactivity(ctx, time.Minute*2)

	requestDone := signal.ExecuteAsync(func() error {
		writer := buf.NewBufferedWriter(conn)
		session.EncodeRequestHeader(request, writer)

		bodyWriter := session.EncodeRequestBody(request, writer)
		firstPayload, err := input.ReadTimeout(time.Millisecond * 500)
		if err != nil && err != buf.ErrReadTimeout {
			return errors.New("failed to get first payload").Base(err).Path("VMess", "Outbound")
		}
		if !firstPayload.IsEmpty() {
			if err := bodyWriter.Write(firstPayload); err != nil {
				return errors.New("failed to write first payload").Base(err).Path("VMess", "Outbound")
			}
			firstPayload.Release()
		}

		if err := writer.SetBuffered(false); err != nil {
			return err
		}

		var inputReader buf.Reader = input
		if request.Command == protocol.RequestCommandTCP {
			inputReader = buf.NewMergingReader(input)
		}

		if err := buf.PipeUntilEOF(timer, inputReader, bodyWriter); err != nil {
			return err
		}

		if request.Option.Has(protocol.RequestOptionChunkStream) {
			if err := bodyWriter.Write(buf.NewLocal(8)); err != nil {
				return err
			}
		}
		return nil
	})

	responseDone := signal.ExecuteAsync(func() error {
		defer output.Close()

		reader := buf.NewBufferedReader(conn)
		header, err := session.DecodeResponseHeader(reader)
		if err != nil {
			return err
		}
		v.handleCommand(rec.Destination(), header.Command)

		conn.SetReusable(header.Option.Has(protocol.ResponseOptionConnectionReuse))

		reader.SetBuffered(false)
		bodyReader := session.DecodeResponseBody(request, reader)
		if err := buf.PipeUntilEOF(timer, bodyReader, output); err != nil {
			return err
		}

		return nil
	})

	if err := signal.ErrorOrFinish2(ctx, requestDone, responseDone); err != nil {
		conn.SetReusable(false)
		return errors.New("connection ends").Base(err).Path("VMess", "Outbound")
	}
	runtime.KeepAlive(timer)

	return nil
}

func init() {
	common.Must(common.RegisterConfig((*Config)(nil), func(ctx context.Context, config interface{}) (interface{}, error) {
		return New(ctx, config.(*Config))
	}))
}
