package transport

var (
	connectionReuse = true
)

// IsConnectionReusable returns true if V2Ray is trying to reuse TCP connections.
func IsConnectionReusable() bool {
	return connectionReuse
}
