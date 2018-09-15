package internet

func isTCPSocket(network string) bool {
	switch network {
	case "tcp", "tcp4", "tcp6":
		return true
	default:
		return false
	}
}

func isUDPSocket(network string) bool {
	switch network {
	case "udp", "udp4", "udp6":
		return true
	default:
		return false
	}
}
