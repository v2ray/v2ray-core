// +build !confonly

package mtproto

import (
	"bytes"
	"context"
	"time"

	"v2ray.com/core"
	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/crypto"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/common/session"
	"v2ray.com/core/common/signal"
	"v2ray.com/core/common/task"
	"v2ray.com/core/features/policy"
	"v2ray.com/core/features/routing"
	"v2ray.com/core/transport/internet"
)

var (
	dcList = []net.Address{
		net.ParseAddress("149.154.175.50"),
		net.ParseAddress("149.154.167.51"),
		net.ParseAddress("149.154.175.100"),
		net.ParseAddress("149.154.167.91"),
		net.ParseAddress("149.154.171.5"),
	}
)

type Server struct {
	user    *protocol.User
	account *Account
	policy  policy.Manager
}

func NewServer(ctx context.Context, config *ServerConfig) (*Server, error) {
	if len(config.User) == 0 {
		return nil, newError("no user configured.")
	}

	user := config.User[0]
	rawAccount, err := config.User[0].GetTypedAccount()
	if err != nil {
		return nil, newError("invalid account").Base(err)
	}
	account, ok := rawAccount.(*Account)
	if !ok {
		return nil, newError("not a MTProto account")
	}

	v := core.MustFromContext(ctx)

	return &Server{
		user:    user,
		account: account,
		policy:  v.GetFeature(policy.ManagerType()).(policy.Manager),
	}, nil
}

func (s *Server) Network() []net.Network {
	return []net.Network{net.Network_TCP}
}

var ctype1 = []byte{0xef, 0xef, 0xef, 0xef}
var ctype2 = []byte{0xee, 0xee, 0xee, 0xee}

func isValidConnectionType(c [4]byte) bool {
	if bytes.Equal(c[:], ctype1) {
		return true
	}
	if bytes.Equal(c[:], ctype2) {
		return true
	}
	return false
}

func (s *Server) Process(ctx context.Context, network net.Network, conn internet.Connection, dispatcher routing.Dispatcher) error {
	sPolicy := s.policy.ForLevel(s.user.Level)

	if err := conn.SetDeadline(time.Now().Add(sPolicy.Timeouts.Handshake)); err != nil {
		newError("failed to set deadline").Base(err).WriteToLog(session.ExportIDToError(ctx))
	}
	auth, err := ReadAuthentication(conn)
	if err != nil {
		return newError("failed to read authentication header").Base(err)
	}
	defer putAuthenticationObject(auth)

	if err := conn.SetDeadline(time.Time{}); err != nil {
		newError("failed to clear deadline").Base(err).WriteToLog(session.ExportIDToError(ctx))
	}

	auth.ApplySecret(s.account.Secret)

	decryptor := crypto.NewAesCTRStream(auth.DecodingKey[:], auth.DecodingNonce[:])
	decryptor.XORKeyStream(auth.Header[:], auth.Header[:])

	ct := auth.ConnectionType()
	if !isValidConnectionType(ct) {
		return newError("invalid connection type: ", ct)
	}

	dcID := auth.DataCenterID()
	if dcID >= uint16(len(dcList)) {
		return newError("invalid datacenter id: ", dcID)
	}

	dest := net.Destination{
		Network: net.Network_TCP,
		Address: dcList[dcID],
		Port:    net.Port(443),
	}

	ctx, cancel := context.WithCancel(ctx)
	timer := signal.CancelAfterInactivity(ctx, cancel, sPolicy.Timeouts.ConnectionIdle)
	ctx = policy.ContextWithBufferPolicy(ctx, sPolicy.Buffer)

	sc := SessionContext{
		ConnectionType: ct,
		DataCenterID:   dcID,
	}
	ctx = ContextWithSessionContext(ctx, sc)

	link, err := dispatcher.Dispatch(ctx, dest)
	if err != nil {
		return newError("failed to dispatch request to: ", dest).Base(err)
	}

	request := func() error {
		defer timer.SetTimeout(sPolicy.Timeouts.DownlinkOnly)

		reader := buf.NewReader(crypto.NewCryptionReader(decryptor, conn))
		return buf.Copy(reader, link.Writer, buf.UpdateActivity(timer))
	}

	response := func() error {
		defer timer.SetTimeout(sPolicy.Timeouts.UplinkOnly)

		encryptor := crypto.NewAesCTRStream(auth.EncodingKey[:], auth.EncodingNonce[:])
		writer := buf.NewWriter(crypto.NewCryptionWriter(encryptor, conn))
		return buf.Copy(link.Reader, writer, buf.UpdateActivity(timer))
	}

	var responseDoneAndCloseWriter = task.OnSuccess(response, task.Close(link.Writer))
	if err := task.Run(ctx, request, responseDoneAndCloseWriter); err != nil {
		common.Interrupt(link.Reader)
		common.Interrupt(link.Writer)
		return newError("connection ends").Base(err)
	}

	return nil
}

func init() {
	common.Must(common.RegisterConfig((*ServerConfig)(nil), func(ctx context.Context, config interface{}) (interface{}, error) {
		return NewServer(ctx, config.(*ServerConfig))
	}))
}
