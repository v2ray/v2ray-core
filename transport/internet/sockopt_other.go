// +build !linux

package internet

func applySocketOptions(fd uintptr, config *SocketConfig) error {
	return nil
}
