// +build !confonly

package trojan

import (
	"context"
	"crypto/tls"
	"io"
	"strconv"
	"time"

	"v2ray.com/core"
	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/errors"
	"v2ray.com/core/common/log"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/protocol"
	udp_proto "v2ray.com/core/common/protocol/udp"
	"v2ray.com/core/common/retry"
	"v2ray.com/core/common/session"
	"v2ray.com/core/common/signal"
	"v2ray.com/core/common/task"
	"v2ray.com/core/features/policy"
	"v2ray.com/core/features/routing"
	"v2ray.com/core/transport/internet"
	"v2ray.com/core/transport/internet/udp"
	"v2ray.com/core/transport/internet/xtls"
)

func init() {
	common.Must(common.RegisterConfig((*ServerConfig)(nil), func(ctx context.Context, config interface{}) (interface{}, error) { // nolint: lll
		return NewServer(ctx, config.(*ServerConfig))
	}))
}

// Server is an inbound connection handler that handles messages in trojan protocol.
type Server struct {
	policyManager policy.Manager
	validator     *Validator
	fallbacks     map[string]map[string]*Fallback // or nil
}

// NewServer creates a new trojan inbound handler.
func NewServer(ctx context.Context, config *ServerConfig) (*Server, error) {
	validator := new(Validator)
	for _, user := range config.Users {
		u, err := user.ToMemoryUser()
		if err != nil {
			return nil, newError("failed to get trojan user").Base(err).AtError()
		}

		if err := validator.Add(u); err != nil {
			return nil, newError("failed to add user").Base(err).AtError()
		}
	}

	v := core.MustFromContext(ctx)
	server := &Server{
		policyManager: v.GetFeature(policy.ManagerType()).(policy.Manager),
		validator:     validator,
	}

	if config.Fallbacks != nil {
		server.fallbacks = make(map[string]map[string]*Fallback)
		for _, fb := range config.Fallbacks {
			if server.fallbacks[fb.Alpn] == nil {
				server.fallbacks[fb.Alpn] = make(map[string]*Fallback)
			}
			server.fallbacks[fb.Alpn][fb.Path] = fb
		}
		if server.fallbacks[""] != nil {
			for alpn, pfb := range server.fallbacks {
				if alpn != "" { // && alpn != "h2" {
					for path, fb := range server.fallbacks[""] {
						if pfb[path] == nil {
							pfb[path] = fb
						}
					}
				}
			}
		}
	}

	return server, nil
}

// Network implements proxy.Inbound.Network().
func (s *Server) Network() []net.Network {
	return []net.Network{net.Network_TCP}
}

// Process implements proxy.Inbound.Process().
func (s *Server) Process(ctx context.Context, network net.Network, conn internet.Connection, dispatcher routing.Dispatcher) error { // nolint: funlen,lll

	sid := session.ExportIDToError(ctx)

	iConn := conn
	if statConn, ok := iConn.(*internet.StatCouterConnection); ok {
		iConn = statConn.Connection
	}

	sessionPolicy := s.policyManager.ForLevel(0)
	if err := conn.SetReadDeadline(time.Now().Add(sessionPolicy.Timeouts.Handshake)); err != nil {
		return newError("unable to set read deadline").Base(err).AtWarning()
	}

	first := buf.New()
	defer first.Release()

	firstLen, err := first.ReadFrom(conn)
	if err != nil {
		return newError("failed to read first request").Base(err)
	}
	newError("firstLen = ", firstLen).AtInfo().WriteToLog(sid)

	bufferedReader := &buf.BufferedReader{
		Reader: buf.NewReader(conn),
		Buffer: buf.MultiBuffer{first},
	}

	var user *protocol.MemoryUser

	apfb := s.fallbacks
	isfb := apfb != nil

	shouldFallback := false
	if firstLen < 58 || first.Byte(56) != '\r' { // nolint: gomnd
		// invalid protocol
		err = newError("not trojan protocol")
		log.Record(&log.AccessMessage{
			From:   conn.RemoteAddr(),
			To:     "",
			Status: log.AccessRejected,
			Reason: err,
		})

		shouldFallback = true
	} else {
		user = s.validator.Get(hexString(first.BytesTo(56))) // nolint: gomnd
		if user == nil {
			// invalid user, let's fallback
			err = newError("not a valid user")
			log.Record(&log.AccessMessage{
				From:   conn.RemoteAddr(),
				To:     "",
				Status: log.AccessRejected,
				Reason: err,
			})

			shouldFallback = true
		}
	}

	if isfb && shouldFallback {
		return s.fallback(ctx, sid, err, sessionPolicy, conn, iConn, apfb, first, firstLen, bufferedReader)
	} else if shouldFallback {
		return newError("invalid protocol or invalid user")
	}

	clientReader := &ConnReader{Reader: bufferedReader}
	if err := clientReader.ParseHeader(); err != nil {
		log.Record(&log.AccessMessage{
			From:   conn.RemoteAddr(),
			To:     "",
			Status: log.AccessRejected,
			Reason: err,
		})
		return newError("failed to create request from: ", conn.RemoteAddr()).Base(err)
	}

	destination := clientReader.Target
	if err := conn.SetReadDeadline(time.Time{}); err != nil {
		return newError("unable to set read deadline").Base(err).AtWarning()
	}

	inbound := session.InboundFromContext(ctx)
	if inbound == nil {
		panic("no inbound metadata")
	}
	inbound.User = user
	sessionPolicy = s.policyManager.ForLevel(user.Level)

	if destination.Network == net.Network_UDP { // handle udp request
		return s.handleUDPPayload(ctx, &PacketReader{Reader: clientReader}, &PacketWriter{Writer: conn}, dispatcher)
	}

	// handle tcp request

	ctx = log.ContextWithAccessMessage(ctx, &log.AccessMessage{
		From:   conn.RemoteAddr(),
		To:     destination,
		Status: log.AccessAccepted,
		Reason: "",
		Email:  user.Email,
	})

	newError("received request for ", destination).WriteToLog(sid)
	return s.handleConnection(ctx, sessionPolicy, destination, clientReader, buf.NewWriter(conn), dispatcher)
}

func (s *Server) handleUDPPayload(ctx context.Context, clientReader *PacketReader, clientWriter *PacketWriter, dispatcher routing.Dispatcher) error { // nolint: lll
	udpServer := udp.NewDispatcher(dispatcher, func(ctx context.Context, packet *udp_proto.Packet) {
		common.Must(clientWriter.WriteMultiBufferWithMetadata(buf.MultiBuffer{packet.Payload}, packet.Source))
	})

	inbound := session.InboundFromContext(ctx)
	user := inbound.User

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			p, err := clientReader.ReadMultiBufferWithMetadata()
			if err != nil {
				if errors.Cause(err) != io.EOF {
					return newError("unexpected EOF").Base(err)
				}
				return nil
			}

			ctx = log.ContextWithAccessMessage(ctx, &log.AccessMessage{
				From:   inbound.Source,
				To:     p.Target,
				Status: log.AccessAccepted,
				Reason: "",
				Email:  user.Email,
			})
			newError("tunnelling request to ", p.Target).WriteToLog(session.ExportIDToError(ctx))

			for _, b := range p.Buffer {
				udpServer.Dispatch(ctx, p.Target, b)
			}
		}
	}
}

func (s *Server) handleConnection(ctx context.Context, sessionPolicy policy.Session,
	destination net.Destination,
	clientReader buf.Reader,
	clientWriter buf.Writer, dispatcher routing.Dispatcher) error {
	ctx, cancel := context.WithCancel(ctx)
	timer := signal.CancelAfterInactivity(ctx, cancel, sessionPolicy.Timeouts.ConnectionIdle)
	ctx = policy.ContextWithBufferPolicy(ctx, sessionPolicy.Buffer)

	link, err := dispatcher.Dispatch(ctx, destination)
	if err != nil {
		return newError("failed to dispatch request to ", destination).Base(err)
	}

	requestDone := func() error {
		defer timer.SetTimeout(sessionPolicy.Timeouts.DownlinkOnly)

		if err := buf.Copy(clientReader, link.Writer, buf.UpdateActivity(timer)); err != nil {
			return newError("failed to transfer request").Base(err)
		}
		return nil
	}

	responseDone := func() error {
		defer timer.SetTimeout(sessionPolicy.Timeouts.UplinkOnly)

		if err := buf.Copy(link.Reader, clientWriter, buf.UpdateActivity(timer)); err != nil {
			return newError("failed to write response").Base(err)
		}
		return nil
	}

	var requestDonePost = task.OnSuccess(requestDone, task.Close(link.Writer))
	if err := task.Run(ctx, requestDonePost, responseDone); err != nil {
		common.Must(common.Interrupt(link.Reader))
		common.Must(common.Interrupt(link.Writer))
		return newError("connection ends").Base(err)
	}

	return nil
}

func (s *Server) fallback(ctx context.Context, sid errors.ExportOption, err error, sessionPolicy policy.Session, connection internet.Connection, iConn internet.Connection, apfb map[string]map[string]*Fallback, first *buf.Buffer, firstLen int64, reader buf.Reader) error { // nolint: lll
	if err := connection.SetReadDeadline(time.Time{}); err != nil {
		newError("unable to set back read deadline").Base(err).AtWarning().WriteToLog(sid)
	}
	newError("fallback starts").Base(err).AtInfo().WriteToLog(sid)

	alpn := ""
	if len(apfb) > 1 || apfb[""] == nil {
		if tlsConn, ok := iConn.(*tls.Conn); ok {
			alpn = tlsConn.ConnectionState().NegotiatedProtocol
			newError("realAlpn = " + alpn).AtInfo().WriteToLog(sid)
		} else if xtlsConn, ok := iConn.(*xtls.Conn); ok {
			alpn = xtlsConn.ConnectionState().NegotiatedProtocol
			newError("realAlpn = " + alpn).AtInfo().WriteToLog(sid)
		}
		if apfb[alpn] == nil {
			alpn = ""
		}
	}
	pfb := apfb[alpn]
	if pfb == nil {
		return newError(`failed to find the default "alpn" config`).AtWarning()
	}

	path := ""
	if len(pfb) > 1 || pfb[""] == nil {
		if firstLen >= 18 && first.Byte(4) != '*' { // not h2c
			firstBytes := first.Bytes()
			for i := 4; i <= 8; i++ { // 5 -> 9
				if firstBytes[i] == '/' && firstBytes[i-1] == ' ' {
					search := len(firstBytes)
					if search > 64 {
						search = 64 // up to about 60
					}
					for j := i + 1; j < search; j++ {
						k := firstBytes[j]
						if k == '\r' || k == '\n' { // avoid logging \r or \n
							break
						}
						if k == ' ' {
							path = string(firstBytes[i:j])
							newError("realPath = " + path).AtInfo().WriteToLog(sid)
							if pfb[path] == nil {
								path = ""
							}
							break
						}
					}
					break
				}
			}
		}
	}
	fb := pfb[path]
	if fb == nil {
		return newError(`failed to find the default "path" config`).AtWarning()
	}

	ctx, cancel := context.WithCancel(ctx)
	timer := signal.CancelAfterInactivity(ctx, cancel, sessionPolicy.Timeouts.ConnectionIdle)
	ctx = policy.ContextWithBufferPolicy(ctx, sessionPolicy.Buffer)

	var conn net.Conn
	if err := retry.ExponentialBackoff(5, 100).On(func() error {
		var dialer net.Dialer
		conn, err = dialer.DialContext(ctx, fb.Type, fb.Dest)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return newError("failed to dial to " + fb.Dest).Base(err).AtWarning()
	}
	defer conn.Close()

	serverReader := buf.NewReader(conn)
	serverWriter := buf.NewWriter(conn)

	postRequest := func() error {
		defer timer.SetTimeout(sessionPolicy.Timeouts.DownlinkOnly)
		if fb.Xver != 0 {
			remoteAddr, remotePort, err := net.SplitHostPort(connection.RemoteAddr().String())
			if err != nil {
				return err
			}
			localAddr, localPort, err := net.SplitHostPort(connection.LocalAddr().String())
			if err != nil {
				return err
			}
			ipv4 := true
			for i := 0; i < len(remoteAddr); i++ {
				if remoteAddr[i] == ':' {
					ipv4 = false
					break
				}
			}
			pro := buf.New()
			defer pro.Release()
			switch fb.Xver {
			case 1:
				if ipv4 {
					common.Must2(pro.Write([]byte("PROXY TCP4 " + remoteAddr + " " + localAddr + " " + remotePort + " " + localPort + "\r\n")))
				} else {
					common.Must2(pro.Write([]byte("PROXY TCP6 " + remoteAddr + " " + localAddr + " " + remotePort + " " + localPort + "\r\n")))
				}
			case 2:
				common.Must2(pro.Write([]byte("\x0D\x0A\x0D\x0A\x00\x0D\x0A\x51\x55\x49\x54\x0A\x21"))) // signature + v2 + PROXY
				if ipv4 {
					common.Must2(pro.Write([]byte("\x11\x00\x0C"))) // AF_INET + STREAM + 12 bytes
					common.Must2(pro.Write(net.ParseIP(remoteAddr).To4()))
					common.Must2(pro.Write(net.ParseIP(localAddr).To4()))
				} else {
					common.Must2(pro.Write([]byte("\x21\x00\x24"))) // AF_INET6 + STREAM + 36 bytes
					common.Must2(pro.Write(net.ParseIP(remoteAddr).To16()))
					common.Must2(pro.Write(net.ParseIP(localAddr).To16()))
				}
				p1, _ := strconv.ParseUint(remotePort, 10, 16)
				p2, _ := strconv.ParseUint(localPort, 10, 16)
				common.Must2(pro.Write([]byte{byte(p1 >> 8), byte(p1), byte(p2 >> 8), byte(p2)}))
			}
			if err := serverWriter.WriteMultiBuffer(buf.MultiBuffer{pro}); err != nil {
				return newError("failed to set PROXY protocol v", fb.Xver).Base(err).AtWarning()
			}
		}
		if err := buf.Copy(reader, serverWriter, buf.UpdateActivity(timer)); err != nil {
			return newError("failed to fallback request payload").Base(err).AtInfo()
		}
		return nil
	}

	writer := buf.NewWriter(connection)

	getResponse := func() error {
		defer timer.SetTimeout(sessionPolicy.Timeouts.UplinkOnly)
		if err := buf.Copy(serverReader, writer, buf.UpdateActivity(timer)); err != nil {
			return newError("failed to deliver response payload").Base(err).AtInfo()
		}
		return nil
	}

	if err := task.Run(ctx, task.OnSuccess(postRequest, task.Close(serverWriter)), task.OnSuccess(getResponse, task.Close(writer))); err != nil {
		common.Must(common.Interrupt(serverReader))
		common.Must(common.Interrupt(serverWriter))
		return newError("fallback ends").Base(err).AtInfo()
	}

	return nil
}
