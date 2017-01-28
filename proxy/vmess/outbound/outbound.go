package outbound

import (
	"context"
	"time"

	"v2ray.com/core/app"
	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/bufio"
	"v2ray.com/core/common/errors"
	"v2ray.com/core/common/log"
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

// VMessOutboundHandler is an outbound connection handler for VMess protocol.
type VMessOutboundHandler struct {
	serverList   *protocol.ServerList
	serverPicker protocol.ServerPicker
}

func New(ctx context.Context, config *Config) (*VMessOutboundHandler, error) {
	space := app.SpaceFromContext(ctx)
	if space == nil {
		return nil, errors.New("VMess|Outbound: No space in context.")
	}

	serverList := protocol.NewServerList()
	for _, rec := range config.Receiver {
		serverList.AddServer(protocol.NewServerSpecFromPB(*rec))
	}
	handler := &VMessOutboundHandler{
		serverList:   serverList,
		serverPicker: protocol.NewRoundRobinServerPicker(serverList),
	}

	return handler, nil
}

// Dispatch implements OutboundHandler.Dispatch().
func (v *VMessOutboundHandler) Process(ctx context.Context, outboundRay ray.OutboundRay) error {
	var rec *protocol.ServerSpec
	var conn internet.Connection

	dialer := proxy.DialerFromContext(ctx)
	err := retry.ExponentialBackoff(5, 100).On(func() error {
		rec = v.serverPicker.PickServer()
		rawConn, err := dialer.Dial(ctx, rec.Destination())
		if err != nil {
			return err
		}
		conn = rawConn

		return nil
	})
	if err != nil {
		log.Warning("VMess|Outbound: Failed to find an available destination:", err)
		return err
	}
	defer conn.Close()

	target := proxy.DestinationFromContext(ctx)
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
		log.Warning("VMess|Outbound: Failed to get user account: ", err)
		return err
	}
	account := rawAccount.(*vmess.InternalAccount)
	request.Security = account.Security

	conn.SetReusable(true)
	if conn.Reusable() { // Conn reuse may be disabled on transportation layer
		request.Option.Set(protocol.RequestOptionConnectionReuse)
	}

	input := outboundRay.OutboundInput()
	output := outboundRay.OutboundOutput()

	session := encoding.NewClientSession(protocol.DefaultIDHash)

	requestDone := signal.ExecuteAsync(func() error {
		writer := bufio.NewWriter(conn)
		session.EncodeRequestHeader(request, writer)

		bodyWriter := session.EncodeRequestBody(request, writer)
		firstPayload, err := input.ReadTimeout(time.Millisecond * 500)
		if err != nil && err != ray.ErrReadTimeout {
			return errors.Base(err).Message("VMess|Outbound: Failed to get first payload.")
		}
		if !firstPayload.IsEmpty() {
			if err := bodyWriter.Write(firstPayload); err != nil {
				return errors.Base(err).Message("VMess|Outbound: Failed to write first payload.")
			}
			firstPayload.Release()
		}

		writer.SetBuffered(false)

		if err := buf.PipeUntilEOF(input, bodyWriter); err != nil {
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

		reader := bufio.NewReader(conn)
		header, err := session.DecodeResponseHeader(reader)
		if err != nil {
			return err
		}
		v.handleCommand(rec.Destination(), header.Command)

		conn.SetReusable(header.Option.Has(protocol.ResponseOptionConnectionReuse))

		reader.SetBuffered(false)
		bodyReader := session.DecodeResponseBody(request, reader)
		if err := buf.Pipe(bodyReader, output); err != nil {
			return err
		}

		return nil
	})

	if err := signal.ErrorOrFinish2(ctx, requestDone, responseDone); err != nil {
		log.Info("VMess|Outbound: Connection ending with ", err)
		conn.SetReusable(false)
		return err
	}

	return nil
}

func init() {
	common.Must(common.RegisterConfig((*Config)(nil), func(ctx context.Context, config interface{}) (interface{}, error) {
		return New(ctx, config.(*Config))
	}))
}
