package ray

import "v2ray.com/core/common/buf"
import "time"

// OutboundRay is a transport interface for outbound connections.
type OutboundRay interface {
	// OutboundInput provides a stream for the input of the outbound connection.
	// The outbound connection shall write all the input until it is closed.
	OutboundInput() InputStream

	// OutboundOutput provides a stream to retrieve the response from the
	// outbound connection. The outbound connection shall close the channel
	// after all responses are receivced and put into the channel.
	OutboundOutput() OutputStream
}

// InboundRay is a transport interface for inbound connections.
type InboundRay interface {
	// InboundInput provides a stream to retrieve the request from client.
	// The inbound connection shall close the channel after the entire request
	// is received and put into the channel.
	InboundInput() OutputStream

	// InboudBound provides a stream of data for the inbound connection to write
	// as response. The inbound connection shall write all the data from the
	// channel until it is closed.
	InboundOutput() InputStream
}

// Ray is an internal tranport channel between inbound and outbound connection.
type Ray interface {
	InboundRay
	OutboundRay
	AddInspector(Inspector)
}

type InputStream interface {
	buf.Reader
	ReadTimeout(time.Duration) (*buf.Buffer, error)
	ForceClose()
}

type OutputStream interface {
	buf.Writer
	Close()
}
