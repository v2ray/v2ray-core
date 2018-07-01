package outbound

//go:generate go run $GOPATH/src/v2ray.com/core/common/errors/errorgen/main.go -pkg outbound -path Proxy,VMess,Outbound

import (
	"context"
	"time"

	"v2ray.com/core/common/session"
	"v2ray.com/core/common/task"

	"v2ray.com/core"
	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/common/retry"
	"v2ray.com/core/common/signal"
	"v2ray.com/core/proxy"
	"v2ray.com/core/proxy/vmess"
	"v2ray.com/core/proxy/vmess/encoding"
	"v2ray.com/core/transport/internet"
)

// Handler is an outbound connection handler for VMess protocol.
type Handler struct {
	serverList   *protocol.ServerList
	serverPicker protocol.ServerPicker
	v            *core.Instance
}

// New creates a new VMess outbound handler.
func New(ctx context.Context, config *Config) (*Handler, error) {
	serverList := protocol.NewServerList()
	for _, rec := range config.Receiver {
		serverList.AddServer(protocol.NewServerSpecFromPB(*rec))
	}
	handler := &Handler{
		serverList:   serverList,
		serverPicker: protocol.NewRoundRobinServerPicker(serverList),
		v:            core.MustFromContext(ctx),
	}

	return handler, nil
}

// Process implements proxy.Outbound.Process().
func (v *Handler) Process(ctx context.Context, link *core.Link, dialer proxy.Dialer) error {
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
		return newError("failed to find an available destination").Base(err).AtWarning()
	}
	defer conn.Close() //nolint: errcheck

	target, ok := proxy.TargetFromContext(ctx)
	if !ok {
		return newError("target not specified").AtError()
	}
	newError("tunneling request to ", target, " via ", rec.Destination()).WriteToLog(session.ExportIDToError(ctx))

	command := protocol.RequestCommandTCP
	if target.Network == net.Network_UDP {
		command = protocol.RequestCommandUDP
	}
	if target.Address.Family().IsDomain() && target.Address.Domain() == "v1.mux.cool" {
		command = protocol.RequestCommandMux
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
		return newError("failed to get user account").Base(err).AtWarning()
	}
	account := rawAccount.(*vmess.InternalAccount)
	request.Security = account.Security

	if request.Security == protocol.SecurityType_AES128_GCM || request.Security == protocol.SecurityType_NONE || request.Security == protocol.SecurityType_CHACHA20_POLY1305 {
		request.Option.Set(protocol.RequestOptionChunkMasking)
	}

	input := link.Reader
	output := link.Writer

	session := encoding.NewClientSession(protocol.DefaultIDHash)
	sessionPolicy := v.v.PolicyManager().ForLevel(request.User.Level)

	ctx, cancel := context.WithCancel(ctx)
	timer := signal.CancelAfterInactivity(ctx, cancel, sessionPolicy.Timeouts.ConnectionIdle)

	requestDone := func() error {
		defer timer.SetTimeout(sessionPolicy.Timeouts.DownlinkOnly)

		writer := buf.NewBufferedWriter(buf.NewWriter(conn))
		if err := session.EncodeRequestHeader(request, writer); err != nil {
			return newError("failed to encode request").Base(err).AtWarning()
		}

		bodyWriter := session.EncodeRequestBody(request, writer)
		if err := buf.CopyOnceTimeout(input, bodyWriter, time.Millisecond*500); err != nil && err != buf.ErrNotTimeoutReader && err != buf.ErrReadTimeout {
			return newError("failed to write first payload").Base(err)
		}

		if err := writer.SetBuffered(false); err != nil {
			return err
		}

		if err := buf.Copy(input, bodyWriter, buf.UpdateActivity(timer)); err != nil {
			return err
		}

		if request.Option.Has(protocol.RequestOptionChunkStream) {
			if err := bodyWriter.WriteMultiBuffer(buf.MultiBuffer{}); err != nil {
				return err
			}
		}

		return nil
	}

	responseDone := func() error {
		defer timer.SetTimeout(sessionPolicy.Timeouts.UplinkOnly)

		reader := &buf.BufferedReader{Reader: buf.NewReader(conn)}
		header, err := session.DecodeResponseHeader(reader)
		if err != nil {
			return newError("failed to read header").Base(err)
		}
		v.handleCommand(rec.Destination(), header.Command)

		reader.Direct = true
		bodyReader := session.DecodeResponseBody(request, reader)

		return buf.Copy(bodyReader, output, buf.UpdateActivity(timer))
	}

	var responseDonePost = task.Single(responseDone, task.OnSuccess(task.Close(output)))
	if err := task.Run(task.WithContext(ctx), task.Parallel(requestDone, responseDonePost))(); err != nil {
		return newError("connection ends").Base(err)
	}

	return nil
}

func init() {
	common.Must(common.RegisterConfig((*Config)(nil), func(ctx context.Context, config interface{}) (interface{}, error) {
		return New(ctx, config.(*Config))
	}))
}
