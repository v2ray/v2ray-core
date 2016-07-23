package encoding

const (
	Version = byte(1)

	AddrTypeIPv4   = byte(0x01)
	AddrTypeIPv6   = byte(0x03)
	AddrTypeDomain = byte(0x02)
)
