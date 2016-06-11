package transport

var (
	connectionReuse = true
)

func IsConnectionReusable() bool {
	return connectionReuse
}
