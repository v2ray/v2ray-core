package transport

var (
	connectionReuse = false
)

func IsConnectionReusable() bool {
	return connectionReuse
}
