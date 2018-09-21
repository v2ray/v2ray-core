package mtproto

import (
	"context"

	"v2ray.com/core"
	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/crypto"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/session"
	"v2ray.com/core/common/task"
	"v2ray.com/core/proxy"
)

type Client struct {
}

func NewClient(ctx context.Context, config *ClientConfig) (*Client, error) {
	return &Client{}, nil
}

func (c *Client) Process(ctx context.Context, link *core.Link, dialer proxy.Dialer) error {
	outbound := session.OutboundFromContext(ctx)
	if outbound == nil || !outbound.Target.IsValid() {
		return newError("unknown destination.")
	}
	dest := outbound.Target
	if dest.Network != net.Network_TCP {
		return newError("not TCP traffic", dest)
	}

	conn, err := dialer.Dial(ctx, dest)
	if err != nil {
		return newError("failed to dial to ", dest).Base(err).AtWarning()
	}
	defer conn.Close() // nolint: errcheck

	auth := NewAuthentication()
	defer putAuthenticationObject(auth)

	request := func() error {
		encryptor := crypto.NewAesCTRStream(auth.EncodingKey[:], auth.EncodingNonce[:])

		var header [HeaderSize]byte
		encryptor.XORKeyStream(header[:], auth.Header[:])
		copy(header[:56], auth.Header[:])

		if _, err := conn.Write(header[:]); err != nil {
			return newError("failed to write auth header").Base(err)
		}

		connWriter := buf.NewWriter(crypto.NewCryptionWriter(encryptor, conn))
		return buf.Copy(link.Reader, connWriter)
	}

	response := func() error {
		decryptor := crypto.NewAesCTRStream(auth.DecodingKey[:], auth.DecodingNonce[:])

		connReader := buf.NewReader(crypto.NewCryptionReader(decryptor, conn))
		return buf.Copy(connReader, link.Writer)
	}

	var responseDoneAndCloseWriter = task.Single(response, task.OnSuccess(task.Close(link.Writer)))
	if err := task.Run(task.WithContext(ctx), task.Parallel(request, responseDoneAndCloseWriter))(); err != nil {
		return newError("connection ends").Base(err)
	}

	return nil
}

func init() {
	common.Must(common.RegisterConfig((*ClientConfig)(nil), func(ctx context.Context, config interface{}) (interface{}, error) {
		return NewClient(ctx, config.(*ClientConfig))
	}))
}
